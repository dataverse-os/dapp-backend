package ceramic

import (
	"bytes"
	"encoding/base64"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ipfs/go-cid"
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
			_, err := BuildDagJWSFromReader(bytes.NewBuffer(buf))
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildDagJWSFromReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestBuildPayloadFromReader(t *testing.T) {
	type args struct {
		base64Encoded string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "common",
			args: args{
				base64Encoded: "omRkYXRhpmdvcHRpb25zeQomZXlKbWIyeGtaWEpPWVcxbElqb2lWVzUwYVhSc1pXUWlMQ0ptYjJ4a1pYSkVaWE5qY21sd2RHbHZiaUk2SWlJc0ltVnVZM0o1Y0hSbFpGTjViVzFsZEhKcFkwdGxlU0k2SW1NM1pUWm1ORGs0T0RVeVpXWXlaV000TVRRd1l6WmhZamN3WkRkallXVTRZbVF3WmpNNU5HUTROemhqTURNNE9UTmpaR1JpTVdVNVlqbGhNekExTVRKaFpHRXlaRFV5TnpjelpUSm1OV001WXpCa01UVTBNak5pTWpOaE1tUXdaV1JrWkdVNVlUY3pZak0wWTJKbE1tTTNNRFJrTWpJMk1UWTBNekZoTkRBM1pEVXpNekl4WkRWaE5EbGpaR1JpTlRjMlkySTJPR05sWVRjd01tRTJObVptTm1abU5HUTRPR05rT1dabFpqTTBaamMxTVROaE56WmtaamhtTnpBNFptWmtNREJpWldJM1lqQXdObUk0WkRZME9UVXdOemxrWlRoaVpUQXpaR1kzTURVMU5qUm1OREJqTTJKallXRTFNV1UyTUdFM01EQmlPR016WW1abFl6RXdNREF3TURBd01EQXdNREF3TURJd05EaGhOalpsWVRkaE5URTVObUk0TW1Fell6VTRZbVJrWkdJMk1UQTBNV0l3WXpZNE16Z3dNREJrWkdZMFpUaGlaalEzTm1FNFpXSXpOV0UxWlRjeU1XUm1NbVpoTXpsak5qRmpORFl5WVRGa056VmxZMk5sTWpBME1tUmpOemcxSWl3aVpHVmpjbmx3ZEdsdmJrTnZibVJwZEdsdmJuTWlPaUpYTTNOcFdUSTVkV1JJU21oWk0xSkNXa2RTZVZwWVRucEphbTlwU1dsM2FXTXpVbWhpYlZKb1kyMVNSR0l5TlRCamJVWnFaRVpTTldOSFZXbFBhVWxwVEVOS2FtRkhSbkJpYVVrMlNXMVdNR0ZIVm5sYVdGWjBTV2wzYVdKWFZqQmhSemxyU1dwdmFVbHBkMmxqUjBaNVdWY3hiR1JIVm5samVVazJWM2xKTm1SWVRteGphMFpyV2toS2JHTXpUV2xZVTNkcFkyMVdNR1JZU25WV2JVWnpaRmRXVlZwWVRqQkphbkEzU1cxT2RtSllRbWhqYlVZd1lqTkphVTlwU1RsSmFYZHBaRzFHYzJSWFZXbFBhVWwzWlVSTmVFMXRWa0pQUkZWNVRucEpNbEpVVGtKUFYxa3lUWHBPUWsxRVRUTk9NazEzV2xkRk5FOUVTWGRQUkZwclRtcFpNazVxV1dsbVdEQnpaWGxLZG1OSFZubFpXRkoyWTJsSk5rbHRSblZhUTBvNVRFaHphVmt5T1hWa1NFcG9XVE5TUWxwSFVubGFXRTU2U1dwdmFVbHBkMmxqTTFKb1ltMVNhR050VWtSaU1qVXdZMjFHYW1SR1VqVmpSMVZwVDJsS1ZGTldaRVpKYVhkcFdUSm9hR0ZYTkdsUGFVcHNaRWRvYkdOdFZqRmlVMGx6U1cweGJHUkhhSFphUTBrMlNXbEpjMGx1UW1oamJVWjBXbGhTYkdOdVRXbFBiSE5wVDI1S2JHTXlPVEZqYlU1c1kzbEtaRXhEU25sYVdGSXhZMjAxVjFsWGVERmFWbEpzWXpOUmFVOXVjMmxaTWpsMFkwZEdlVmxZVW5aamFVazJTVzFPZG1KdVVtaGhWelY2U1dsM2FXUnRSbk5rVjFWcFQybEthbHBZU21oaVYyeHFUMms0ZGt0cU9YUmlNbEpzWWtReGNtRnVjSE5PYldneVdtNUthV1I2V21wT1dFWnJaVzVrY0U5WFZucGxTRm93VFZoWk1XSllVakJPTWpsclRqSm9hVTFxYXpCT2VsbDVUa2N4ZFU1SVZYZGpiVEY0VFZoS2IwOVhSblZoYlU1MVpVaG5hV1pZTUhObGVVcDJZMGRXZVZsWVVuWmphVWsyU1cxR2RWcERTamxNU0hOcFdUSTVkV1JJU21oWk0xSkNXa2RTZVZwWVRucEphbTlwU1dsM2FXTXpVbWhpYlZKb1kyMVNSR0l5TlRCamJVWnFaRVpTTldOSFZXbFBhVXBVVTFaa1JrbHBkMmxaTW1ob1lWYzBhVTlwU214a1IyaHNZMjFXTVdKVFNYTkpiVEZzWkVkb2RscERTVFpKYVVselNXNUNhR050Um5SYVdGSnNZMjVOYVU5c2MybFBia3BzWXpJNU1XTnRUbXhqZVVwa1RFTktlVnBZVWpGamJUVlhXVmQ0TVZwV1VteGpNMUZwVDI1emFWa3lPWFJqUjBaNVdWaFNkbU5wU1RaSmJVNTJZbTVTYUdGWE5YcEphWGRwWkcxR2MyUlhWV2xQYVVwcVdsaEthR0pYYkdwUGFUaDJTMm81ZEdJeVVteGlSREZ5WVc1d2MwNXRhREphYmtwcFpIcGFhazV0Um10T00yeHJZbXBDYjJGVVVqSmtSMFowWlVSS01rNXFTWGRoUjFKdVpGUmFjMkpJUlRCUFYyZDVUMGhLYlZwRVdtcGpla0Y1V25wT2FtSlhORFZsYlVWcFpsZ3djMlY1U25aalIxWjVXVmhTZG1OcFNUWkpiVVoxV2tOS09VeEljMmxaTWpsMVpFaEthRmt6VWtKYVIxSjVXbGhPZWtscWIybEphWGRwWXpOU2FHSnRVbWhqYlZKRVlqSTFNR050Um1wa1JsSTFZMGRWYVU5cFNsUlRWbVJHU1dsM2FWa3lhR2hoVnpScFQybEtiR1JIYUd4amJWWXhZbE5KYzBsdE1XeGtSMmgyV2tOSk5rbHBTWE5KYmtKb1kyMUdkRnBZVW14amJrMXBUMnh6YVU5dVNteGpNamt4WTIxT2JHTjVTbVJNUTBwNVdsaFNNV050TlZkWlYzZ3hXbFpTYkdNelVXbFBibk5wV1RJNWRHTkhSbmxaV0ZKMlkybEpOa2x0VG5aaWJsSm9ZVmMxZWtscGQybGtiVVp6WkZkVmFVOXBTbXBhV0Vwb1lsZHNhazlwT0haTGFqbDBZakpTYkdKRU1YSmhibkJ6VG0xb01scHVTbWxrZWxwcVRucFplbVJYU210aFJ6a3paVzFHZGsxSE1EQmxXRUUwVGtkT05HVnRTbTFpYlhodlRrZG9hMkZVVm1oaVNFWjJUa2hzZVZwWFNuUlpla0o0WTBkd2EyRlVWV2xtV0RGa0lpd2laVzVqY25sd2RHVmtJam9pUVVSWWNXYzNXbGRLVkVwNVNYUllaakZpTTJ0T01YcFliVFJCTVZoQlJVdzRXRXRUZVRWdFdYRnRaMlJ1UzI5RVprTmlWM2RoVFhkNWQwWlRlVWRtTmxaVlZGaHJlWEZ6ZW1sNFVrSkhMV3d6T0d4cFRWVmZOREp3WldJME0yWjBZa2RTYjJOblRXRk1NblZMTWpadFVYSXpNVlJDY0dSWlFWbDRSbGxWVnpVaWZRaWNyZWF0ZWRBdHgYMjAyMy0wNi0xMVQxMToyNDoyMy4wMjlaaXVwZGF0ZWRBdHgYMjAyMy0wNi0xMVQxMToyNDoyMy4wMjlaamFwcFZlcnNpb25lMC4yLjBqZm9sZGVyVHlwZQFwY29udGVudEZvbGRlcklkc4FkdGVtcGZoZWFkZXKjZW1vZGVsWCjOAQIBhQESIBzfRtp3sMC9cDyAKBAj8zq/N2Xd5IJPzEEoVKyvyJM1ZnVuaXF1ZUwrbXEamTAna5uAMqJrY29udHJvbGxlcnOBeDtkaWQ6cGtoOmVpcDE1NToxOjB4MzEyZUE4NTI3MjZFM0E5ZjYzM0EwMzc3YzBlYTg4MjA4NmQ2NjY2Ng==",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf, _ := base64.StdEncoding.DecodeString(tt.args.base64Encoded)
			BuildPayloadFromReader(bytes.NewBuffer(buf))
		})
	}
}

func TestConvertFrom(t *testing.T) {
	type args struct {
		base64Encoded string
	}
	tests := []struct {
		name       string
		args       args
		wantCommit Commit
		wantErr    bool
	}{
		{
			name: "common",
			args: args{
				base64Encoded: "omdwYXlsb2FkWCQBcRIgGtqDmUTOFzN1JfGKAS9bP1KGrprc6FzbjWjnGQv9yjlqc2lnbmF0dXJlc4GiaXByb3RlY3RlZFjMeyJhbGciOiJFZERTQSIsImNhcCI6ImlwZnM6Ly9iYWZ5cmVpZm5qdno3NDZ1azJrcWhpdHpqN2I1b3hhNmRjcDR3N2J1ZzJhbnhqbHNlYXhiM2R0bG9qeSIsImtpZCI6ImRpZDprZXk6ejZNa3UxOUdnY2gxTlliUmlOREZQVGFwcmQ0aUZKVlVncThTTnpNNnNSTEZESk45I3o2TWt1MTlHZ2NoMU5ZYlJpTkRGUFRhcHJkNGlGSlZVZ3E4U056TTZzUkxGREpOOSJ9aXNpZ25hdHVyZVhAoy7KX9UhdAIpYq2cvN8Kjv+Ka3ps6GzjYHt6afsu8qYYI7pS4k9/Ln013MBMnUTe7eM0/+YuAKnnOmsO8a5ECw==",
			},
			wantCommit: Commit{
				Link:    cid.MustParse("bafyreia23kbzsrgoc4zxkjprrias6wz7kkdk5gw45bonxdli44mqx7okhe"),
				Payload: cid.MustParse("bafyreia23kbzsrgoc4zxkjprrias6wz7kkdk5gw45bonxdli44mqx7okhe"),
				Signatures: []Signature{
					{
						Header:    nil,
						Protected: []byte(`{"alg":"EdDSA","cap":"ipfs://bafyreifnjvz746uk2kqhitzj7b5oxa6dcp4w7bug2anxjlseaxb3dtlojy","kid":"did:key:z6Mku19Ggch1NYbRiNDFPTaprd4iFJVUgq8SNzM6sRLFDJN9#z6Mku19Ggch1NYbRiNDFPTaprd4iFJVUgq8SNzM6sRLFDJN9"}`),
						Signature: hexutil.MustDecode("0xa32eca5fd52174022962ad9cbcdf0a8eff8a6b7a6ce86ce3607b7a69fb2ef2a61823ba52e24f7f2e7d35dcc04c9d44deede334ffe62e00a9e73a6b0ef1ae440b"),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf, _ := base64.StdEncoding.DecodeString(tt.args.base64Encoded)
			dagJws, _ := BuildDagJWSFromReader(bytes.NewBuffer(buf))
			gotCommit, err := ConvertFrom(dagJws)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCommit, tt.wantCommit) {
				t.Errorf("ConvertFrom() = %v\nwant %v", gotCommit, tt.wantCommit)
			}
		})
	}
}
