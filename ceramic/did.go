package ceramic

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"

	"github.com/multiformats/go-multibase"
	"github.com/multiformats/go-multicodec"
)

type DefaultDIDGenerator struct{}

func (DefaultDIDGenerator) GenerateDID(key string) (did string, err error) {
	return GenerateDID(key)
}

var didEncoder = multibase.MustNewEncoder(multibase.Base58BTC)

func GenerateDID(key string) (did string, err error) {
	var seed []byte
	if seed, err = hex.DecodeString(key); err != nil {
		return
	}
	ed25519Key := ed25519.NewKeyFromSeed(seed)
	pub := ed25519Key.Public().(ed25519.PublicKey)
	var buf bytes.Buffer
	buf.Write([]byte{byte(multicodec.Ed25519Pub), 0x01})
	buf.Write(pub)
	did = fmt.Sprintf("did:key:%s", didEncoder.Encode(buf.Bytes()))
	return
}
