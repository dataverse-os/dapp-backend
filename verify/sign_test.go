package verify_test

import (
	"crypto/ecdsa"
	"testing"

	"github.com/dataverse-os/dapp-backend/verify"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/samber/lo"
)

func TestCheckSign(t *testing.T) {
	type args struct {
		data []byte
		sig  string
		key  *ecdsa.PublicKey
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "common",
			args: args{
				data: []byte("I want create a dataverse Dapp:\nName:dTwit\nCeramic Url:https://ceramic.dtwit.com\n"),
				sig:  "0x462efe78ae51f25a8516d3ce14a84eb7366d191ae95a8c3f067f34750d73b68f4174b29c9c94ed0ca2c0008b437fbed5030b311b83f0e39ebe5baa4a7d2053861b",
				key:  &lo.Must(crypto.HexToECDSA("c0eb4f72a47364ac981c9db9636f4809108401126146c3319f2c27286d453b90")).PublicKey,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := verify.CheckSign(tt.args.data, tt.args.sig, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("CheckSign() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckSignWithHexPublicKey(t *testing.T) {
	type args struct {
		origin string
		signed string
	}
	tests := []struct {
		name       string
		args       args
		wantHexKey string
		wantErr    bool
	}{
		{
			name: "common",
			args: args{
				origin: "I want create a dataverse Dapp:\nName:dTwit\nCeramic Url:https://ceramic.dtwit.com\n",
				signed: "0x462efe78ae51f25a8516d3ce14a84eb7366d191ae95a8c3f067f34750d73b68f4174b29c9c94ed0ca2c0008b437fbed5030b311b83f0e39ebe5baa4a7d2053861b",
			},
			wantHexKey: crypto.PubkeyToAddress(lo.Must(crypto.HexToECDSA("c0eb4f72a47364ac981c9db9636f4809108401126146c3319f2c27286d453b90")).PublicKey).Hex(),
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHexKey, err := verify.ExportPublicKey(tt.args.origin, tt.args.signed)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckSignWithHexPublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHexKey != tt.wantHexKey {
				t.Errorf("CheckSignWithHexPublicKey() = %v, want %v", gotHexKey, tt.wantHexKey)
			}
		})
	}
}

func TestSignData(t *testing.T) {
	type args struct {
		data       []byte
		privateKey *ecdsa.PrivateKey
	}
	tests := []struct {
		name          string
		args          args
		wantSignature string
		wantErr       bool
	}{
		{
			name: "common",
			args: args{
				data:       []byte("I want create a dataverse Dapp:\nName:dTwit\nCeramic Url:https://ceramic.dtwit.com\n"),
				privateKey: lo.Must(crypto.HexToECDSA("c0eb4f72a47364ac981c9db9636f4809108401126146c3319f2c27286d453b90")),
			},
			wantSignature: "0x462efe78ae51f25a8516d3ce14a84eb7366d191ae95a8c3f067f34750d73b68f4174b29c9c94ed0ca2c0008b437fbed5030b311b83f0e39ebe5baa4a7d2053861b",
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSignature, err := verify.SignData(tt.args.data, tt.args.privateKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSignature != tt.wantSignature {
				t.Errorf("SignData() = %v, want %v", gotSignature, tt.wantSignature)
			}
		})
	}
}
