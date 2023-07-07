package ceramic_test

import (
	"context"
	"crypto/ecdsa"
	"testing"

	"github.com/dataverse-os/dapp-backend/ceramic"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/samber/lo"
)

func TestCheckIfCeramicAdminKey(t *testing.T) {
	type args struct {
		ctx    context.Context
		config []byte
		key    *ecdsa.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "common",
			args: args{
				ctx:    context.Background(),
				config: []byte(`{"http-api":{"admin-dids":["did:key:z6MkrHGM2j8Jtx1eP5NMaNgiG9oEz4YDShnZYA7k1BfBq7H6"]}}`),
				key:    lo.Must(crypto.HexToECDSA("c0eb4f72a47364ac981c9db9636f4809108401126146c3319f2c27286d453b90")),
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ceramic.CheckIfCeramicAdminKey(tt.args.ctx, tt.args.config, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckIfCeramicAdminKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckIfCeramicAdminKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrivateKeyToDID(t *testing.T) {
	type args struct {
		ctx context.Context
		key *ecdsa.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		wantDid string
		wantErr bool
	}{
		{
			name: "common",
			args: args{
				ctx: context.Background(),
				key: lo.Must(crypto.HexToECDSA("c0eb4f72a47364ac981c9db9636f4809108401126146c3319f2c27286d453b90")),
			},
			wantDid: "did:key:z6MkrHGM2j8Jtx1eP5NMaNgiG9oEz4YDShnZYA7k1BfBq7H6",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDid, err := ceramic.PrivateKeyToDID(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrivateKeyToDID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDid != tt.wantDid {
				t.Errorf("PrivateKeyToDID() = %v, want %v", gotDid, tt.wantDid)
			}
		})
	}
}
