package tests_test

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestVerify(t *testing.T) {
	if os.Getenv("DAPP_BACKEND") == "" {
		t.Skip("skip case without dapp-backend daemon")
	}
	privateKey, _ := crypto.HexToECDSA("c0eb4f72a47364ac981c9db9636f4809108401126146c3319f2c27286d453b90")
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	t.Run("sth", func(t *testing.T) {
		fmt.Println(crypto.PubkeyToAddress(*publicKeyECDSA))
	})

	var msg string
	var sig string
	t.Run("模拟前端生成签名", func(t *testing.T) {
		data := []byte("I want create a dataverse Dapp:\n" + "Name:dTwit\n" + "Ceramic Url:https://ceramic.dtwit.com\n")
		signature, _ := crypto.Sign(crypto.Keccak256Hash(data).Bytes(), privateKey)
		msg = string(data)
		sig = hexutil.Encode(signature)
		fmt.Println(msg)
		fmt.Println(sig)

		req, _ := http.NewRequest(http.MethodPost, os.Getenv("DAPP_BACKEND")+"/dataverse/dapp", bytes.NewBufferString(string(data)))
		req.Header.Set("dataverse-sig", sig)
		req.Header.Set("dataverse-nonce", "123456")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(body))
	})

	t.Run("还原公钥address", func(t *testing.T) {
		// 前端传来string类型的msg sig后，
		data := []byte(msg)
		signature := hexutil.MustDecode(sig)
		hash := crypto.Keccak256Hash(data)
		sigPublicKeyECDSA, _ := crypto.SigToPub(hash.Bytes(), signature)
		addr := crypto.PubkeyToAddress(*sigPublicKeyECDSA)
		fmt.Println("还原出来的address ： ", addr)
		// 这个函数要在我们的dapp-table验证，还要dapp-backend中验证
		// 1. 在我们的dapp-table验证, 然后切割msg中的ceramil url字段， 发送一个请求给那个url （这个url就是他们的dapp-backend后端）
		// 2. 他们的dapp-backend调用现在的函数验证成功后返回给我们的dapp-table
		// 3. 现在就可以保存一条dapp数据库记录。 addr 为这个dapp的公钥。

	})

	t.Run("验证：", func(t *testing.T) {
		// 前端传来string类型的msg sig后，
		data := []byte(msg)
		signature := hexutil.MustDecode(sig)
		hash := crypto.Keccak256Hash(data)
		signatureNoRecoverID := signature[:len(signature)-1] // remove recovery id
		verified := crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)
		fmt.Println("验证传入的msg和sig是不是由本地.privatekey文件中的私钥签名的: ", verified)
	})

	// dapp backend在验证后，用自己的.privatekey再签名一个成功消息还给dapp table，
	t.Run("dapp backend在验证成功后重新签名：", func(t *testing.T) {
		// 上面模拟前端签名的函数
	})

	// dapp-table收到响应值（签名）后也要验证一下
}
