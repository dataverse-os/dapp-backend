package ceramic

import (
	"crypto/sha256"
	"encoding/base64"
)

func MessageHash(msg []byte) string {
	hash := sha256.New()
	hash.Write(msg)
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}
