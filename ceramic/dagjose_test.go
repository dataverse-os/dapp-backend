package ceramic

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func TestBuildDagJWSFromReader(t *testing.T) {
	type args struct {
		base64Encoded string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "common",
			args: args{
				base64Encoded: "omdwYXlsb2FkWCQBcRIgGtqDmUTOFzN1JfGKAS9bP1KGrprc6FzbjWjnGQv9yjlqc2lnbmF0dXJlc4GiaXByb3RlY3RlZFjMeyJhbGciOiJFZERTQSIsImNhcCI6ImlwZnM6Ly9iYWZ5cmVpZm5qdno3NDZ1azJrcWhpdHpqN2I1b3hhNmRjcDR3N2J1ZzJhbnhqbHNlYXhiM2R0bG9qeSIsImtpZCI6ImRpZDprZXk6ejZNa3UxOUdnY2gxTlliUmlOREZQVGFwcmQ0aUZKVlVncThTTnpNNnNSTEZESk45I3o2TWt1MTlHZ2NoMU5ZYlJpTkRGUFRhcHJkNGlGSlZVZ3E4U056TTZzUkxGREpOOSJ9aXNpZ25hdHVyZVhAoy7KX9UhdAIpYq2cvN8Kjv+Ka3ps6GzjYHt6afsu8qYYI7pS4k9/Ln013MBMnUTe7eM0/+YuAKnnOmsO8a5ECw==",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf, _ := base64.StdEncoding.DecodeString(tt.args.base64Encoded)
			_, err := DecodeDagJWSFromReader(bytes.NewBuffer(buf))
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildDagJWSFromReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
