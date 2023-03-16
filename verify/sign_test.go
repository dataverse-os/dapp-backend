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
				data: []byte("I want create a dataverse Dapp:\n" + "Name:dTwit\n" + "Ceramic Url:https://ceramic.dtwit.com\n"),
				sig:  "0x4af1babce1fe5a0a547d7a3d019393ba26efbe2a9248bf6353fd7b51e75086f16b45c3e713539632741e29c846b492b6732f9102af84b31c3df34d1d4f16be9f00",
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
