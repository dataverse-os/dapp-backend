package dapp_test

import (
	"crypto/ecdsa"
	"math/rand"
	"testing"

	"github.com/dataverse-os/dapp-backend/internal/dapp"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/samber/lo"
)

func TestLookupAndUpdateUserModelVersion(t *testing.T) {
	var (
		pubKey        = lo.Must(crypto.GenerateKey()).PublicKey
		tempModelName = "temp"
		randomVersion = rand.Uint64()
	)

	dapp.InitBolt()
	defer dapp.BoltDB.Close()

	t.Run("lookup model version (not exist)", func(t *testing.T) {
		gotVersion, err := dapp.LookupUserModelVersion(&pubKey, tempModelName)
		if err != nil {
			t.Fatal(err)
		}
		if gotVersion != -1 {
			t.Fatalf("failed to lookup a non-exist model version, got: %d want: -1", gotVersion)
		}
	})
	t.Run("update model version", func(t *testing.T) {
		err := dapp.UpdateUserModelVersion(&pubKey, tempModelName, randomVersion)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("lookup model version (exist)", func(t *testing.T) {
		gotVersion, err := dapp.LookupUserModelVersion(&pubKey, tempModelName)
		if err != nil {
			t.Fatal(err)
		}
		if gotVersion != int64(randomVersion) {
			t.Fatalf("failed to lookup a non-exist model version, got: %d want: %d", gotVersion, randomVersion)
		}
	})
}

func TestLookupUserModelVersion(t *testing.T) {
	type args struct {
		pubKey    *ecdsa.PublicKey
		modelName string
	}
	tests := []struct {
		name        string
		args        args
		wantVersion int64
		wantErr     bool
	}{
		{
			name: "common",
			args: args{
				pubKey:    &lo.Must(crypto.GenerateKey()).PublicKey,
				modelName: "tempModelName",
			},
			wantVersion: -1,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dapp.InitBolt()
			defer dapp.BoltDB.Close()
			gotVersion, err := dapp.LookupUserModelVersion(tt.args.pubKey, tt.args.modelName)
			if (err != nil) != tt.wantErr {
				t.Errorf("LookupUserModelVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVersion != tt.wantVersion {
				t.Errorf("LookupUserModelVersion() = %v, want %v", gotVersion, tt.wantVersion)
			}
		})
	}
}

func TestUpdateUserModelVersion(t *testing.T) {
	type args struct {
		pubKey    *ecdsa.PublicKey
		modelName string
		version   uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "common",
			args: args{
				pubKey:    &lo.Must(crypto.GenerateKey()).PublicKey,
				modelName: "tempModelName",
				version:   1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dapp.InitBolt()
			defer dapp.BoltDB.Close()
			if err := dapp.UpdateUserModelVersion(tt.args.pubKey, tt.args.modelName, tt.args.version); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserModelVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
