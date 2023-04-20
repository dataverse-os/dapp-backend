package ceramic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

var _ ceramicClient = (*CliBinding)(nil)

type CliBinding struct{}

func (CliBinding) CreateComposite(ctx context.Context, schema string, ceramic string, key string) (composite string, err error) {
	var (
		graphqlPath string
		output      []byte
	)
	if graphqlPath, _, err = NewFile([]byte(schema)); err != nil {
		return
	}
	cmd := exec.CommandContext(ctx, "composedb", "composite:create",
		graphqlPath,
	)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("DID_PRIVATE_KEY=%s", key),
		fmt.Sprintf("CERAMIC_URL=%s", ceramic),
	)
	if output, err = cmd.CombinedOutput(); err != nil {
		err = fmt.Errorf("output: %s\n err: %s", string(output), err)
		return
	}
	if !bytes.Contains(output, []byte("✔ Creating the composite... Done!")) {
		err = fmt.Errorf("fail to create composite")
		log.Println(string(output))
		return
	}
	arr := bytes.Split(output, []byte("\n"))
	if !json.Valid(arr[len(arr)-2]) {
		err = fmt.Errorf("fail to create composite, not valid json")
	}
	composite = string(arr[len(arr)-2])
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
