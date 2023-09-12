package ceramic

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/crypto"
)

func NewSession(ceramic, key string) (ses Session, err error) {
	ses = Session{
		URLString:      ceramic,
		AdminKeyString: key,
	}
	if ses.URL, err = url.Parse(ceramic); err != nil {
		err = fmt.Errorf("failed to parse ceramic url '%s' as url with error: %s", ceramic, err)
		return
	}

	if ses.AdminKey, err = crypto.HexToECDSA(key); err != nil {
		err = fmt.Errorf("failed to parse ceramic admin key with error: %s", err)
		return
	}
	return
}

type Session struct {
	URL            *url.URL
	URLString      string
	AdminKey       *ecdsa.PrivateKey
	AdminKeyString string
}

var Default = &NodeJSBinding{}

func DeployStreamModel(ctx context.Context, schema string, sess Session) (composite string, streamID string, err error) {
	if composite, err = Default.CreateComposite(ctx, schema, sess); err != nil {
		return
	}
	if streamID, err = ExtarctStreamID(composite); err != nil {
		return
	}
	return
}

type DIDGenerator interface {
	GenerateDID(key string) (did string, err error)
}

type ClientInterface interface {
	CreateComposite(ctx context.Context, schema string, sess Session) (composite string, err error)
	CheckSyntax(ctx context.Context, schema string) (err error)
	CheckAdminAccess(ctx context.Context, sess Session) (err error)
	GetIndexedModels(ctx context.Context, sess Session) (streamIDs []string, err error)
}
