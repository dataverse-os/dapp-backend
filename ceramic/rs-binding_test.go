package ceramic

import (
	"testing"
)

func TestRustBinding_GenerateDID(t *testing.T) {
	type fields struct {
		ClientInterface ClientInterface
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantDid string
		wantErr bool
	}{
		{
			name: "common",
			args: args{
				key: "61c5ed6b2a619e21d7d0d0a9b9a591e4c0f014c3f25eb1d26c1b53332f96afe5",
			},
			wantDid: "did:key:z6MkjSnks3PuMFQhJHS6NfwD3tHfkx6sSGxHjzAQhN113rZj",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RustBinding{
				ClientInterface: tt.fields.ClientInterface,
			}
			gotDid, err := r.GenerateDID(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("RustBinding.GenerateDID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDid != tt.wantDid {
				t.Errorf("RustBinding.GenerateDID() = %v, want %v", gotDid, tt.wantDid)
			}
		})
	}
}
