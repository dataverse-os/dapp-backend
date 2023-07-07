package ceramic

import (
	"context"
)

var Default = &NodeJSBinding{}

func GenerateComposite(ctx context.Context, schema string, url string, key string) (composite string, streamID string, err error) {
	if composite, err = Default.CreateComposite(ctx, schema, url, key); err != nil {
		return
	}
	if streamID, err = ExtarctStreamID(composite); err != nil {
		return
	}
	return
}

type ClientInterface interface {
	CreateComposite(ctx context.Context, schema string, ceramic string, key string) (composite string, err error)
	CheckSyntax(ctx context.Context, schema string) (err error)
	GenerateDID(ctx context.Context, key string) (did string, err error)
	CheckAdminAccess(ctx context.Context, ceramic string, key string) (err error)
}
