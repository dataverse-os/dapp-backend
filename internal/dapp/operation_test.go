package dapp

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
)

func TestDeployStreamModels(t *testing.T) {
	if os.Getenv("CERAMIC_URL") == "" || os.Getenv("CERAMIC_ADMIN_KEY") == "" {
		t.Skip("skip case without ceramic secret")
	}
	type args struct {
		ctx        context.Context
		id         uuid.UUID
		schemas    []StreamModel
		ceramicURL string
		key        string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "common",
			args: args{
				ctx: context.Background(),
				id:  uuid.New(),
				schemas: []StreamModel{
					{
						Schema: `type tempModelName @createModel(accountRelation: LIST, description: "TempModel") {
  num: Int!
  name: String! @string(maxLength: 100)
}`,
						IsPublicDomain: false,
						Encryptable:    []string{"name"},
					},
				},
				ceramicURL: os.Getenv("CERAMIC_URL"),
				key:        os.Getenv("CERAMIC_ADMIN_KEY"),
			},
			wantErr: false,
		},
		{
			name: "missing required directive",
			args: args{
				ctx: context.Background(),
				id:  uuid.New(),
				schemas: []StreamModel{
					{
						Schema: `type tempModelName @createModel(accountRelation: LIST, description: "TempModel") {
  num: Int!
  name: String!
}`,
						IsPublicDomain: false,
						Encryptable:    []string{"name"},
					},
				},
				ceramicURL: os.Getenv("CERAMIC_URL"),
				key:        os.Getenv("CERAMIC_ADMIN_KEY"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DeployStreamModels(tt.args.ctx, tt.args.id, tt.args.schemas, tt.args.ceramicURL, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeployStreamModels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
