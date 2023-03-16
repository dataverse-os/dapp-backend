package dapp

import (
	"bytes"
	"context"
	"dapp-backend/ceramic"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

func GenerateCompositeJsonWithGraphql(ctx context.Context, schema []byte) (streamID string, modelPath string, err error) {
	var (
		graphqlPath string
		output      []byte
	)
	if graphqlPath, _, err = NewFile(schema); err != nil {
		return
	}
	modelFile, err := os.CreateTemp("", "")
	if err != nil {
		return
	}

	modelPath = modelFile.Name()
	if output, err = ceramic.Warper(ctx, "composedb", "composite:create", graphqlPath, "--output="+modelPath); err != nil {
		return
	}
	if bytes.Contains(output, []byte("GraphQLError [Object]: Syntax Error")) {
		err = ErrorGraphQLSyntaxError
		return
	}
	var modelBody []byte
	if modelBody, err = io.ReadAll(modelFile); err != nil {
		return
	}
	defer modelFile.Close()
	for k := range gjson.GetBytes(modelBody, "models").Map() {
		streamID = k
	}
	if streamID == "" {
		err = fmt.Errorf("error generate stream")
	}
	return
}

func DeployCompositeJson(ctx context.Context, modelPath string) (err error) {
	var output []byte
	if output, err = ceramic.Warper(ctx, "composedb", "composite:deploy", modelPath); err != nil {
		return
	}
	arr := strings.Split(string(output), "\n")
	if !strings.Contains(string(output), "Deploying the composite... Done!") || len(arr) != 5 {
		err = fmt.Errorf("create and deploy composite fail")
		return
	}
	return
}

func CreateStreamWithGraphql(ctx context.Context, schema []byte) (streamID string, err error) {
	var modelPath string
	if streamID, modelPath, err = GenerateCompositeJsonWithGraphql(ctx, schema); err != nil {
		return
	}
	if err = DeployCompositeJson(ctx, modelPath); err != nil {
		return
	}
	return
}

func NewFile(content []byte) (path string, file *os.File, err error) {
	file, err = os.CreateTemp("", "")
	if err != nil {
		return
	}
	path = file.Name()
	defer file.Close()
	if _, err = file.Write(content); err != nil {
		return
	}
	return
}
