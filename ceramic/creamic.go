package ceramic

import (
	"bufio"
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tidwall/gjson"
)

func Warper(ctx context.Context, name string, args ...string) (output []byte, err error) {
	cmd := exec.CommandContext(ctx, name, args...)
	if output, err = cmd.CombinedOutput(); err != nil {
		err = fmt.Errorf("output: %s\n err: %s", string(output), err)
		return
	}
	return
}

func PrivateKeyToDID(ctx context.Context, key *ecdsa.PrivateKey) (did string, err error) {
	var output []byte
	if output, err = Warper(ctx, "composedb", "did:from-private-key", hex.EncodeToString(crypto.FromECDSA(key))); err != nil {
		return
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(output))
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "did:") {
			did = scanner.Text()
			return
		}
	}
	err = fmt.Errorf("cannot extract did from composedb cli output")
	return
}

var CeramicConfigPath = os.Getenv("CERAMIC_CONFIG_PATH")

func CheckIfCeramicAdminKeyFromFile(ctx context.Context, key *ecdsa.PrivateKey) (_ bool, err error) {
	var configBytes []byte
	if configBytes, err = os.ReadFile(filepath.Join(CeramicConfigPath, ".ceramic/daemon.config.json")); err != nil {
		return
	}
	return CheckIfCeramicAdminKey(ctx, configBytes, key)
}

func CheckIfCeramicAdminKey(ctx context.Context, config []byte, key *ecdsa.PrivateKey) (_ bool, err error) {
	var did string
	if did, err = PrivateKeyToDID(ctx, key); err != nil {
		return
	}
	adminKeys := gjson.GetBytes(config, "http-api.admin-dids")
	for _, v := range adminKeys.Array() {
		if v.String() == did {
			return true, nil
		}
	}
	return
}
