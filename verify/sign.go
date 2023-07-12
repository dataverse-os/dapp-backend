package verify

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
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

func ExportPublicKeyHex(origin string, signed string) (hexKey string, err error) {
	signature := hexutil.MustDecode(signed)
	signature[64] -= 27
	keyFromSig, err := crypto.SigToPub(accounts.TextHash([]byte(origin)), signature)
	if err != nil {
		return
	}
	hexKey = crypto.PubkeyToAddress(*keyFromSig).Hex()
	return
}

func ExportPublicKey(origin string, signed string) (keyFromSig *ecdsa.PublicKey, err error) {
	signature := hexutil.MustDecode(signed)
	signature[64] -= 27
	keyFromSig, err = crypto.SigToPub(accounts.TextHash([]byte(origin)), signature)
	if err != nil {
		return
	}
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
