package dapp

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/dataverse-os/dapp-backend/ceramic"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

var (
	maxWorkers = func() int {
		m := os.Getenv("MAX_NODEJS_WORKERS")
		if m == "" {
			return runtime.GOMAXPROCS(0) * 4
		} else {
			return lo.Must(strconv.Atoi(m))
		}
	}()
	sem = make(chan struct{}, maxWorkers)
)

func DeployStreamModels(ctx context.Context, id uuid.UUID, schemas []StreamModel, sess ceramic.Session) (result []ModelResult, err error) {
	result = make([]ModelResult, len(schemas))
	eg, _ := errgroup.WithContext(ctx)
	for i := range schemas {
		if err = CheckEncryptable(schemas[i]); err != nil {
			return
		}
		// check encryptable and add encrypted field
		if !schemas[i].IsPublicDomain && len(schemas[i].Encryptable) != 0 {
			if schemas[i].Schema, err =
				ceramic.AddCustomField([]byte(schemas[i].Schema), &EncryptedField); err != nil {
				return
			}
		}
		// add dataverse dapp id to model description prefix
		if schemas[i].Schema, err =
			ceramic.SchemaModifyFn([]byte(schemas[i].Schema), ceramic.OriginModifyFn, func(old string) string {
				return fmt.Sprintf("Dataverse: %s | %s", id, old)
			}); err != nil {
			return
		}
		result[i].Schema = schemas[i].Schema

		// parallel deploy model
		schemaIndex := i
		eg.Go(func() error {

			sem <- struct{}{}
			defer func() { <-sem }()

			// deploy modified model to ceramic node
			if _, result[schemaIndex].StreamID, err =
				ceramic.DeployStreamModel(ctx, schemas[schemaIndex].Schema, sess); err != nil {
				return err
			}
			return nil
		})
	}

	if err = eg.Wait(); err != nil {
		return
	}
	return
}
