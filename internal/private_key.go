package internal

import (
	"crypto/ecdsa"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
)

var (
	EnvPrivateKey = os.Getenv("PRIVATE_KEY")
	PrivateKey    *ecdsa.PrivateKey
)

func init() {
	var err error
	if EnvPrivateKey != "" {
		if PrivateKey, err = crypto.HexToECDSA(EnvPrivateKey); err != nil {
			log.Fatalln(err)
		}
	}
	if PrivateKey, err = crypto.GenerateKey(); err != nil {
		log.Fatalln(err)
	}
}
