package ceramic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	jsscripts "github.com/dataverse-os/dapp-backend/js-scripts"
)

var _ ClientInterface = (*NodeJSBinding)(nil)

type NodeJSBinding struct{}

var (
	tempCheckSyntaxScript   *os.File
	tempDeployModelScript   *os.File
	tempAdminAccessScript   *os.File
	tempIndexedModelsScript *os.File
)

func init() {
	var err error
	if tempCheckSyntaxScript, err = initFile(jsscripts.CheckSyntax); err != nil {
		log.Panicln(err)
	}
	if tempDeployModelScript, err = initFile(jsscripts.DeployModel); err != nil {
		log.Panicln(err)
	}
	if tempAdminAccessScript, err = initFile(jsscripts.AdminAccess); err != nil {
		log.Panicln(err)
	}
	if tempIndexedModelsScript, err = initFile(jsscripts.IndexedModels); err != nil {
		log.Panicln(err)
	}
}

func initFile(content string) (file *os.File, err error) {
	if file, err = os.CreateTemp("", ""); err != nil {
		return
	}
	defer file.Close()
	if _, err = file.WriteString(content); err != nil {
		return
	}
	return
}

func (*NodeJSBinding) CheckSyntax(ctx context.Context, schema string) (err error) {
	var buffer bytes.Buffer
	if err = json.NewEncoder(&buffer).Encode(schema); err != nil {
		return
	}
	cmd := exec.CommandContext(ctx, "node", tempCheckSyntaxScript.Name(), buffer.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	if len(out) != 0 {
		err = fmt.Errorf("check syntax with error from ceramic: %s", strings.TrimSuffix(string(out), "\n"))
	}
	return
}

func (*NodeJSBinding) CreateComposite(ctx context.Context, schema string, sess Session) (composite string, err error) {
	data := map[string]any{
		"schema":  schema,
		"ceramic": sess.URLString,
		"key":     sess.AdminKeyString,
	}
	var buffer bytes.Buffer
	if err = json.NewEncoder(&buffer).Encode(data); err != nil {
		return
	}
	cmd := exec.CommandContext(ctx, "node", tempDeployModelScript.Name(), buffer.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(out))
		return
	}
	composite = strings.TrimSuffix(string(out), "\n")
	return
}

func (*NodeJSBinding) CheckAdminAccess(ctx context.Context, sess Session) (err error) {
	data := map[string]any{
		"ceramic": sess.URLString,
		"key":     sess.AdminKeyString,
	}
	var buffer bytes.Buffer
	if err = json.NewEncoder(&buffer).Encode(data); err != nil {
		return
	}
	cmd := exec.CommandContext(ctx, "node", tempAdminAccessScript.Name(), buffer.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	if bytes.HasPrefix(out, []byte("Error")) {
		err = errors.New(strings.TrimSuffix(string(out), "\n"))
	}
	return
}

func (*NodeJSBinding) GetIndexedModels(ctx context.Context, sess Session) (streamIDs []string, err error) {
	data := map[string]any{
		"ceramic": sess.URLString,
		"key":     sess.AdminKeyString,
	}
	var buffer bytes.Buffer
	if err = json.NewEncoder(&buffer).Encode(data); err != nil {
		return
	}
	cmd := exec.CommandContext(ctx, "node", tempIndexedModelsScript.Name(), buffer.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	if bytes.HasPrefix(out, []byte("Error")) {
		err = errors.New(strings.TrimSuffix(string(out), "\n"))
		return
	}
	if err = json.NewDecoder(bytes.NewBuffer(out)).Decode(&streamIDs); err != nil {
		return
	}
	return
}
