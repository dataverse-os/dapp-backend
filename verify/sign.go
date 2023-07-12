package verify

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func CheckSign(data []byte, sig string, key *ecdsa.PublicKey) error {
	signature := hexutil.MustDecode(sig)
	signature[64] -= 27
	signatureNoRecoverID := signature[:len(signature)-1] // remove recovery id
	publicKeyBytes := crypto.FromECDSAPub(key)
	verified := crypto.VerifySignature(publicKeyBytes, accounts.TextHash(data), signatureNoRecoverID)
	if !verified {
		return fmt.Errorf("error not verified")
	}
	return nil
}

func ExportPublicKey(origin []byte, signatureString string) (keyFromSig *ecdsa.PublicKey, err error) {
	signature := hexutil.MustDecode(signatureString)
	signature[64] -= 27
	keyFromSig, err = crypto.SigToPub(accounts.TextHash(origin), signature)
	if err != nil {
		return
	}
	return
}

func ExportAddress(origin []byte, signatureString string) (address common.Address, err error) {
	keyFromSig, err := ExportPublicKey(origin, signatureString)
	if err != nil {
		return
	}
	address = crypto.PubkeyToAddress(*keyFromSig)
	return
}

func ExportAddressHex(origin []byte, signed string) (addressHex string, err error) {
	address, err := ExportAddress(origin, signed)
	if err != nil {
		return
	}
	addressHex = address.Hex()
	return
}

func SignData(data []byte, privateKey *ecdsa.PrivateKey) (hexSig string, err error) {
	var signature []byte
	if signature, err = crypto.Sign(accounts.TextHash(data), privateKey); err != nil {
		return
	}
	signature[64] += 27
	hexSig = hexutil.Encode(signature)
	return
}
