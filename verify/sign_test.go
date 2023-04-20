package verify_test

import (
	"crypto/ecdsa"
	"dapp-backend/verify"
	"testing"

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
				sig:  "0x4af1babce1fe5a0a547d7a3d019393ba26efbe2a9248bf6353fd7b51e75086f16b45c3e713539632741e29c846b492b6732f9102af84b31c3df34d1d4f16be9f00",
				key:  &lo.Must(crypto.HexToECDSA("c0eb4f72a47364ac981c9db9636f4809108401126146c3319f2c27286d453b90")).PublicKey,
			},
			wantErr: false,
		},
		{
			name: "common",
			args: args{
				data: []byte("hello"),
				sig:  "0xb32c89d7f2d9ec26b7394d8d81a367416feb9ec1d4a1387cc3ed0f465c4178d451a581c231bf901a7cfce8844a0a227aa1fabdc9686a546f8f000b3dab1937411c",
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
				signed: "0x4af1babce1fe5a0a547d7a3d019393ba26efbe2a9248bf6353fd7b51e75086f16b45c3e713539632741e29c846b492b6732f9102af84b31c3df34d1d4f16be9f00",
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
