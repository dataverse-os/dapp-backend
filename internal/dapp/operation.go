package dapp

import (
	"context"
	"fmt"

	"github.com/dataverse-os/dapp-backend/ceramic"
	"github.com/google/uuid"
)

func DeployStreamModels(ctx context.Context, id uuid.UUID, schemas []StreamModel, sess ceramic.Session) (result []ModelResult, err error) {
	result = make([]ModelResult, len(schemas))
	for i := range schemas {
		if err = CheckEncryptable(schemas[i]); err != nil {
			return
		}
		// check encryptable and add encrypted field
		if !schemas[i].IsPublicDomain && len(schemas[i].Encryptable) != 0 {
			if schemas[i].Schema, err = ceramic.AddCustomField([]byte(schemas[i].Schema), &EncryptedField); err != nil {
				return
			}
		}
		// add dataverse dapp id to model description prefix
		if schemas[i].Schema, err = ceramic.SchemaModifyFn([]byte(schemas[i].Schema), ceramic.OriginModifyFn, func(old string) string {
			return fmt.Sprintf("Dataverse: %s | %s", id, old)
		}); err != nil {
			return
		}
		// deploy modified model to ceramic node
		if _, result[i].StreamID, err = ceramic.DeployStreamModel(context.Background(), schemas[i].Schema, sess); err != nil {
			return
		}
		result[i].Schema = schemas[i].Schema
	}
	return
}
