package ceramic

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/samber/lo"
)

func TestNodeJSBinding_CheckSyntax(t *testing.T) {
	type args struct {
		ctx    context.Context
		schema string
	}
	tests := []struct {
		name    string
		n       *NodeJSBinding
		args    args
		wantErr bool
	}{
		{
			name: "common",
			args: args{
				ctx: context.Background(),
				schema: `type contentFolders @createModel(accountRelation: LIST, description: "ContentFolder") {
	author: DID! @documentAccount
	version: CommitID! @documentVersion
	indexFolderId: String! @string(maxLength: 1000)
	mirrors: String! @string(maxLength: 300000000)
  }`,
			},
			wantErr: false,
		},
		{
			name: "common error",
			args: args{
				ctx: context.Background(),
				schema: `type contentFolders @createModel() {
	author: DID! @documentAccount
  }`,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NodeJSBinding{}
			if err := n.CheckSyntax(tt.args.ctx, tt.args.schema); (err != nil) != tt.wantErr {
				t.Errorf("NodeJSBinding.CheckSyntax() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNodeJSBinding_CheckAdminAccess(t *testing.T) {
	if os.Getenv("CERAMIC_URL") == "" || os.Getenv("CERAMIC_ADMIN_KEY") == "" {
		t.Skip("skip case without ceramic secret")
	}
	type args struct {
		ctx  context.Context
		sess Session
	}
	tests := []struct {
		name    string
		n       *NodeJSBinding
		args    args
		wantErr bool
	}{
		{
			name: "common",
			args: args{
				ctx: context.Background(),
				sess: Session{
					URLString:      os.Getenv("CERAMIC_URL"),
					AdminKeyString: os.Getenv("CERAMIC_ADMIN_KEY"),
				},
			},
			wantErr: false,
		},
		{
			name: "common error",
			args: args{
				ctx: context.Background(),
				sess: Session{
					URLString: os.Getenv("CERAMIC_URL"),
					// random generated key
					AdminKeyString: hex.EncodeToString(crypto.FromECDSA(lo.Must(crypto.GenerateKey()))),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NodeJSBinding{}
			if err := n.CheckAdminAccess(tt.args.ctx, tt.args.sess); (err != nil) != tt.wantErr {
				t.Errorf("NodeJSBinding.CheckAdminAccess() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNodeJSBinding_CreateComposite(t *testing.T) {
	if os.Getenv("CERAMIC_URL") == "" || os.Getenv("CERAMIC_ADMIN_KEY") == "" {
		t.Skip("skip case without ceramic secret")
	}
	type args struct {
		ctx    context.Context
		schema string
		sess   Session
	}
	tests := []struct {
		name    string
		n       *NodeJSBinding
		args    args
		wantErr bool
	}{
		{
			name: "common",
			n:    &NodeJSBinding{},
			args: args{
				ctx: context.Background(),
				schema: `type testSchema1 @createModel(accountRelation: LIST, description: "ContentFolder") {
	author: DID! @documentAccount
  }`,
				sess: Session{
					URLString:      os.Getenv("CERAMIC_URL"),
					AdminKeyString: os.Getenv("CERAMIC_ADMIN_KEY"),
				},
			},
			wantErr: false,
		},
		{
			name: "common error",
			n:    &NodeJSBinding{},
			args: args{
				ctx: context.Background(),
				schema: `type testSchema2 @111createModel(accountRelation: LIST, description: "ContentFolder") {
	author: DID! @documentAccount
  }`,
				sess: Session{
					URLString:      os.Getenv("CERAMIC_URL"),
					AdminKeyString: os.Getenv("CERAMIC_ADMIN_KEY"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NodeJSBinding{}
			gotComposite, err := n.CreateComposite(tt.args.ctx, tt.args.schema, tt.args.sess)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeJSBinding.CreateComposite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !json.Valid([]byte(gotComposite)) {
				t.Errorf("NodeJSBinding.CreateComposite() error, got: %v", gotComposite)
			}
		})
	}
}

func TestNodeJSBinding_GetIndexedModels(t *testing.T) {
	if os.Getenv("CERAMIC_URL") == "" || os.Getenv("CERAMIC_ADMIN_KEY") == "" {
		t.Skip("skip case without ceramic secret")
	}
	type args struct {
		ctx  context.Context
		sess Session
	}
	tests := []struct {
		name    string
		n       *NodeJSBinding
		args    args
		wantErr bool
	}{
		{
			name: "common",
			n:    &NodeJSBinding{},
			args: args{
				ctx: context.Background(),
				sess: Session{
					URLString:      os.Getenv("CERAMIC_URL"),
					AdminKeyString: os.Getenv("CERAMIC_ADMIN_KEY"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NodeJSBinding{}
			gotStreamIDs, err := n.GetIndexedModels(tt.args.ctx, tt.args.sess)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeJSBinding.GetIndexedModels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotStreamIDs) == 0 {
				t.Errorf("NodeJSBinding.GetIndexedModels() got empty indexed models, got: %s", gotStreamIDs)
				return
			}
		})
	}
}
