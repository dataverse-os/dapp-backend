package ceramic

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
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

func TestNodeJSBinding_GenerateDID(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		n       *NodeJSBinding
		args    args
		wantDid string
		wantErr bool
	}{
		{
			name: "common",
			args: args{
				ctx: context.Background(),
				key: "61c5ed6b2a619e21d7d0d0a9b9a591e4c0f014c3f25eb1d26c1b53332f96afe5",
			},
			wantDid: "did:key:z6MkjSnks3PuMFQhJHS6NfwD3tHfkx6sSGxHjzAQhN113rZj",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NodeJSBinding{}
			gotDid, err := n.GenerateDID(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeJSBinding.GenerateDID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !strings.HasPrefix(gotDid, "did:") {
				t.Errorf("NodeJSBinding.GenerateDID() invalid did: %v", gotDid)
			}
		})
	}
}

func TestNodeJSBinding_CheckAdminAccess(t *testing.T) {
	if os.Getenv("CERAMIC_URL") == "" || os.Getenv("CERAMIC_ADMIN_KEY") == "" {
		t.Skip("skip case without ceramic secret")
	}
	type args struct {
		ctx     context.Context
		ceramic string
		key     string
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
				ctx:     context.Background(),
				ceramic: os.Getenv("CERAMIC_URL"),
				key:     os.Getenv("CERAMIC_ADMIN_KEY"),
			},
			wantErr: false,
		},
		{
			name: "common error",
			args: args{
				ctx:     context.Background(),
				ceramic: os.Getenv("CERAMIC_URL"),
				// random generated key
				key: hex.EncodeToString(crypto.FromECDSA(lo.Must(crypto.GenerateKey()))),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NodeJSBinding{}
			if err := n.CheckAdminAccess(tt.args.ctx, tt.args.ceramic, tt.args.key); (err != nil) != tt.wantErr {
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
		ctx     context.Context
		schema  string
		ceramic string
		key     string
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
				ceramic: os.Getenv("CERAMIC_URL"),
				key:     os.Getenv("CERAMIC_ADMIN_KEY"),
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
				ceramic: os.Getenv("CERAMIC_URL"),
				key:     os.Getenv("CERAMIC_ADMIN_KEY"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NodeJSBinding{}
			gotComposite, err := n.CreateComposite(tt.args.ctx, tt.args.schema, tt.args.ceramic, tt.args.key)
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
