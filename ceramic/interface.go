package ceramic

import (
	"context"
	"fmt"
	"os"

	"github.com/tidwall/gjson"
)

var CeramicClient ceramicClient = &CliBinding{}

func GenerateComposite(ctx context.Context, schema string) (composite string, streamID string, err error) {
	if composite, err = CeramicClient.CreateComposite(ctx, schema, os.Getenv("CERAMIC_URL"), os.Getenv("DID_PRIVATE_KEY")); err != nil {
		return
	}
	if streamID, err = ExtarctStreamID(composite); err != nil {
		return
	}
	return
}

type ceramicClient interface {
	CreateComposite(ctx context.Context, schema string, ceramic string, key string) (composite string, err error)
}

func ExtarctStreamID(composite string) (streamID string, err error) {
	for k := range gjson.Get(composite, "models").Map() {
		streamID = k
		return
	}
	err = fmt.Errorf("streamID not found")
	return
}
