package dapp

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/google/uuid"
)

var (
	ErrorGraphQLSyntaxError = fmt.Errorf("GraphQLError [Object]: Syntax Error")
)

var (
	descriptionRe = regexp.MustCompile("(?i)@createModel\\(accountRelation: LIST, description: \"((?:[^\"]|\"\")*)\"\\)")
	nameRe        = regexp.MustCompile(`(?i)type\s+([A-Za-z0-9_]+)\s+@createModel`)
)

func AddDappNameInDescription(schema []byte, name string) (result []byte, err error) {
	res := descriptionRe.FindSubmatch(schema)
	if len(res) != 2 {
		err = fmt.Errorf("schema syntax error, please check your description in @createModel")
		return
	}
	newCreateModel := bytes.Replace(res[0], res[1], []byte(fmt.Sprintf("Dataverse:%s | %s", name, res[1])), 1)
	result = bytes.Replace(schema, res[0], newCreateModel, 1)
	return
}

func AppendAppIDInDescription(schema []byte, appID uuid.UUID) (result []byte, err error) {
	res := descriptionRe.FindSubmatch(schema)
	if len(res) != 2 {
		err = fmt.Errorf("schema syntax error, please check your description in @createModel")
		return
	}
	newCreateModel := bytes.Replace(res[0], res[1], []byte(fmt.Sprintf("%s | dataverse:%s", res[1], appID.String())), 1)
	result = bytes.Replace(schema, res[0], newCreateModel, 1)
	return
}

func ExtractModelName(schema []byte) (name string, err error) {
	res := nameRe.FindSubmatch(schema)
	if len(res) != 2 {
		err = fmt.Errorf("schema syntax error, please check your model name")
		return
	}
	name = string(res[1])
	return
}

func ReplaceModelName(schema []byte, name string) (result []byte, err error) {
	res := nameRe.FindSubmatch(schema)
	if len(res) != 2 {
		err = fmt.Errorf("schema syntax error, please check your model name")
		return
	}
	newModelName := bytes.Replace(res[0], res[1], []byte(name), 1)
	result = bytes.Replace(schema, res[0], newModelName, 1)
	return
}

func ReplaceModelNameFunc(schema []byte, fn func(modelName string) string) (result []byte, err error) {
	var modelName string
	if modelName, err = ExtractModelName(schema); err != nil {
		return
	}
	return ReplaceModelName(schema, fn(modelName))
}
