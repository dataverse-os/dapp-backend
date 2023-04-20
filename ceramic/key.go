package ceramic

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateKey() {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalln(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("[Private Key]    >>> %s\n", hex.EncodeToString(crypto.FromECDSA(privateKey)))
	fmt.Printf("[Public Address] >>> %s\n\n", fromAddress.Hex())
	fmt.Printf("If you are using docker, set -e DID_PRIVATE_KEY={PRIVATE_KEY_HERE} flag\n")
}
