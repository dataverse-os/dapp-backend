package verify_test

import (
	"dapp-backend/verify"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestParseBodyMiddleware(t *testing.T) {
	type args struct {
		ctx  *gin.Context
		keys map[string]any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "common",
			args: args{
				ctx: &gin.Context{
					Request: &http.Request{
						Body: io.NopCloser(strings.NewReader(`I want create a dataverse Dapp:
Name:dTwit
Ceramic Url:https://ceramic.dtwit.com`)),
					},
				},
				keys: map[string]any{
					"DATAVERSE_I_WANT_CREATE_A_DATAVERSE_DAPP": "",
					"DATAVERSE_CERAMIC_URL":                    "https://ceramic.dtwit.com",
					"DATAVERSE_NAME":                           "dTwit",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verify.ParseBodyMiddleware(tt.args.ctx)
			if !reflect.DeepEqual(tt.args.ctx.Keys, tt.args.keys) {
				t.Errorf("keys extract from body = %v, want %v", tt.args.ctx.Keys, tt.args.keys)
			}
		})
	}
}
