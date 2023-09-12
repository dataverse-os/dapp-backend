package ceramic

import (
	"context"

	"github.com/dataverse-os/dapp-backend/lib"
)

var _ ClientInterface = (*RustBinding)(nil)

type RustBinding struct {
	ClientInterface
}

func (*RustBinding) CheckAdminAccess(ctx context.Context, sess Session) error {
	_, err := lib.GetCeramicNodeStatus(sess.URLString, sess.AdminKeyString)
	if err != nil {
		return err
	}
	return nil
}

func (*RustBinding) GetIndexedModels(ctx context.Context, sess Session) (streamIDs []string, err error) {
	panic("unimplemented")
}

func (*RustBinding) GenerateDID(key string) (did string, err error) {
	return lib.GenerateDID(key)
}
