package verify

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

func CheckMiddleware(key *ecdsa.PrivateKey) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			data bytes.Buffer
			err  error
		)
		if _, err = io.Copy(&data, ctx.Request.Body); err != nil {
			ctx.AbortWithError(400, err)
		}
		if err = CheckSign(data.Bytes(), ctx.GetHeader("dataverse-sig"), &key.PublicKey); err != nil {
			ctx.AbortWithError(403, err)
		}
	}
}

func CheckSign(data []byte, sig string, key *ecdsa.PublicKey) error {
	signature := hexutil.MustDecode(sig)
	hash := crypto.Keccak256Hash(data)
	signatureNoRecoverID := signature[:len(signature)-1] // remove recovery id
	publicKeyBytes := crypto.FromECDSAPub(key)
	verified := crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)
	if !verified {
		return fmt.Errorf("error not verified")
	}
	return nil
}

func SignData(data []byte, privateKey *ecdsa.PrivateKey) (signature []byte, err error) {
	return crypto.Sign(crypto.Keccak256Hash(data).Bytes(), privateKey)
}
