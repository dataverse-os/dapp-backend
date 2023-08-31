package dapptable

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

type getDappByModelIDQuery struct {
	GetDapp Dapp `graphql:"getDapp(modelId: $modelId)"`
}

type Dapp struct {
	Id     string
	Models []struct {
		ModelName string
		Streams   []struct {
			ModelId     string
			Encryptable []string
			Version     int
		}
	}
}

func GetDappByModelID(ctx context.Context, modelId string) (dapp Dapp, err error) {
	client := graphql.NewClient("https://gateway.dataverse.art/v1/dapp-table/graphql", nil)
	var query getDappByModelIDQuery
	if err = client.Query(ctx, &query, map[string]any{
		"modelId": modelId,
	}); err != nil {
		return
	}
	dapp = query.GetDapp
	return
}
