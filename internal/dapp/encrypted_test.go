package dapp

import (
	"testing"
)

func TestCheckEncryptable(t *testing.T) {
	type args struct {
		schema StreamModel
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "common",
			args: args{
				schema: StreamModel{
					Schema: `type temp @createModel(accountRelation: LIST, description: "empty") {
  field1: String!
}`,
					IsPublicDomain: false,
					Encryptable: []string{
						"field1",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "common no encryptable",
			args: args{
				schema: StreamModel{
					Schema: `type temp @createModel(accountRelation: LIST, description: "empty") {
  field1: String!
}`,
					IsPublicDomain: false,
					Encryptable:    []string{},
				},
			},
			wantErr: false,
		},
		{
			name: "public domain but encryptable",
			args: args{
				schema: StreamModel{
					Schema: `type temp @createModel(accountRelation: LIST, description: "empty") {
  field1: String!
}`,
					IsPublicDomain: true,
					Encryptable: []string{
						"field1",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "encryptable but with encrypted field",
			args: args{
				schema: StreamModel{
					Schema: `type temp @createModel(accountRelation: LIST, description: "empty") {
  field1: String!
  encrypted: String!
}`,
					IsPublicDomain: true,
					Encryptable: []string{
						"field1",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "encryptable but with encrypted field",
			args: args{
				schema: StreamModel{
					Schema: `type temp @createModel(accountRelation: LIST, description: "empty") {
  field1: DID! @documentAccount
}`,
					IsPublicDomain: true,
					Encryptable: []string{
						"field1",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckEncryptable(tt.args.schema); (err != nil) != tt.wantErr {
				t.Errorf("CheckEncryptable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
