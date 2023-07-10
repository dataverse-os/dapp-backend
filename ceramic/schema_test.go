package ceramic_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/dataverse-os/dapp-backend/ceramic"
)

func Test_ExtractModelName(t *testing.T) {
	type args struct {
		schema []byte
	}
	tests := []struct {
		name     string
		args     args
		wantName string
		wantErr  bool
	}{
		{
			name: "tempModelName",
			args: args{
				schema: []byte(`type tempModelName @createModel(accountRelation: LIST) {}`),
			},
			wantName: "tempModelName",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, err := ceramic.ExtractModelName(tt.args.schema)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractModelName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotName != tt.wantName {
				t.Errorf("ExtractModelName() = %v, want %v", gotName, tt.wantName)
			}
		})
	}
}

func TestModelNameSyntaxCheck(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "common",
			args: args{
				name: "abc123",
			},
			wantErr: false,
		},
		{
			name: "number prefix",
			args: args{
				name: "123abc",
			},
			wantErr: true,
		},
		{
			name: "dash prefix",
			args: args{
				name: "_abc123",
			},
			wantErr: true,
		},
		{
			name: "dash",
			args: args{
				name: "a-bc123",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ceramic.ModelNameSyntaxCheck(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("ModelNameSyntaxCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSchemaModifyFn(t *testing.T) {
	type args struct {
		body    []byte
		nameMod func(old string) string
		despMod func(old string) string
	}
	tests := []struct {
		name       string
		wantResult string
		args       args
		wantErr    bool
	}{
		{
			name: "common",
			args: args{
				body: []byte(`type temp @createModel(accountRelation: SINGLE, description: "desp here") {}`),
				nameMod: func(old string) string {
					return fmt.Sprintf("prefix_%s", old)
				},
				despMod: func(old string) string {
					return fmt.Sprintf("prefix | %s", old)
				},
			},
			wantResult: `type prefix_temp @createModel(accountRelation: SINGLE, description: "prefix | desp here") {}`,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := ceramic.SchemaModifyFn(tt.args.body, tt.args.nameMod, tt.args.despMod)
			if (err != nil) != tt.wantErr {
				t.Errorf("SchemaModifyFn() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("SchemaModifyFn() = %v, want %v", string(gotResult), string(tt.wantResult))
			}
		})
	}
}

func TestFormatSchema(t *testing.T) {
	type args struct {
		src []byte
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
		wantErr    bool
	}{
		{
			name: "common",
			args: args{
				src: []byte(`type multiLineAttr
				@createModel(accountRelation: LIST, description: "desp here") {
				multiLineField: [String!]!
				  @list(minLength: 1, maxLength: 10000)
				  @string(maxLength: 2000)
				deleted: Boolean
			  }`),
			},
			wantResult: `type multiLineAttr @createModel(accountRelation: LIST, description: "desp here") {
  multiLineField: [String!]! @list(minLength: 1, maxLength: 10000) @string(maxLength: 2000)
  deleted: Boolean
}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := ceramic.FormatSchema(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("FormatSchema() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
