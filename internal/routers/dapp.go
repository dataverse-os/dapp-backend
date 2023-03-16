package routers

import (
	"context"
	"dapp-backend/dapp"
	"fmt"
)

var defaultModels = [][]byte{
	dapp.ContentFloder, dapp.IndexFile, dapp.IndexFolder,
}

func CreateDappModels(ctx context.Context, name string) (schemas map[string]string, err error) {
	schemas = make(map[string]string)
	for _, schema := range defaultModels {
		if schema, err = dapp.ReplaceModelNameFunc(schema, func(modelName string) string {
			return fmt.Sprintf("%s_%s", name, modelName)
		}); err != nil {
			return
		}
		var modelPath, streamID string
		if streamID, modelPath, err = dapp.GenerateCompositeJsonWithGraphql(ctx, schema); err != nil {
			return
		}
		if err = dapp.DeployCompositeJson(ctx, modelPath); err != nil {
			return
		}
		schemas[streamID] = string(schema)
	}
	return
}
