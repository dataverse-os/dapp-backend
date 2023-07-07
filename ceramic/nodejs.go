package ceramic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"

	jsscripts "github.com/dataverse-os/dapp-backend/js-scripts"
)

var _ ClientInterface = (*NodeJSBinding)(nil)

type NodeJSBinding struct{}

var (
	tempCheckSyntaxScript *os.File
	tempDeployModelScript *os.File
	tempGenerateDIDScript *os.File
	tempAdminAccessScript *os.File
)

func init() {
	var err error
	if tempCheckSyntaxScript, err = initFile(jsscripts.CheckSyntax); err != nil {
		log.Panicln(err)
	}
	if tempDeployModelScript, err = initFile(jsscripts.DeployModel); err != nil {
		log.Panicln(err)
	}
	if tempGenerateDIDScript, err = initFile(jsscripts.GenerateDID); err != nil {
		log.Panicln(err)
	}
	if tempAdminAccessScript, err = initFile(jsscripts.AdminAccess); err != nil {
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
		err = errors.New(strings.TrimSuffix(string(out), "\n"))
	}
	return
}

func (*NodeJSBinding) CreateComposite(ctx context.Context, schema string, ceramic string, key string) (composite string, err error) {
	data := map[string]any{
		"schema":  schema,
		"ceramic": ceramic,
		"key":     key,
	}
	var buffer bytes.Buffer
	if err = json.NewEncoder(&buffer).Encode(data); err != nil {
		return
	}
	cmd := exec.CommandContext(ctx, "node", tempDeployModelScript.Name(), buffer.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	composite = strings.TrimSuffix(string(out), "\n")
	return
}

func (*NodeJSBinding) GenerateDID(ctx context.Context, key string) (did string, err error) {
	cmd := exec.CommandContext(ctx, "node", tempGenerateDIDScript.Name(), key)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	did = strings.TrimSuffix(string(out), "\n")
	return
}

func (*NodeJSBinding) CheckAdminAccess(ctx context.Context, ceramic string, key string) (err error) {
	data := map[string]any{
		"ceramic": ceramic,
		"key":     key,
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
