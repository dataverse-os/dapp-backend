package ceramic_test

import (
	"bytes"
	"encoding/base64"
	"reflect"
	"testing"

	"github.com/dataverse-os/dapp-backend/ceramic"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ipfs/go-cid"
	"github.com/samber/lo"
)

func TestGenesisCommit_DecodeFromBlockReader(t *testing.T) {
	type args struct {
		base64Encoded string
	}
	tests := []struct {
		name       string
		wantCommit *ceramic.GenesisCommitPayload
		args       args
		wantErr    bool
	}{
		{
			name: "common",
			args: args{base64Encoded: "omRkYXRhpmdvcHRpb25zeQomZXlKbWIyeGtaWEpPWVcxbElqb2lWVzUwYVhSc1pXUWlMQ0ptYjJ4a1pYSkVaWE5qY21sd2RHbHZiaUk2SWlJc0ltVnVZM0o1Y0hSbFpGTjViVzFsZEhKcFkwdGxlU0k2SW1NM1pUWm1ORGs0T0RVeVpXWXlaV000TVRRd1l6WmhZamN3WkRkallXVTRZbVF3WmpNNU5HUTROemhqTURNNE9UTmpaR1JpTVdVNVlqbGhNekExTVRKaFpHRXlaRFV5TnpjelpUSm1OV001WXpCa01UVTBNak5pTWpOaE1tUXdaV1JrWkdVNVlUY3pZak0wWTJKbE1tTTNNRFJrTWpJMk1UWTBNekZoTkRBM1pEVXpNekl4WkRWaE5EbGpaR1JpTlRjMlkySTJPR05sWVRjd01tRTJObVptTm1abU5HUTRPR05rT1dabFpqTTBaamMxTVROaE56WmtaamhtTnpBNFptWmtNREJpWldJM1lqQXdObUk0WkRZME9UVXdOemxrWlRoaVpUQXpaR1kzTURVMU5qUm1OREJqTTJKallXRTFNV1UyTUdFM01EQmlPR016WW1abFl6RXdNREF3TURBd01EQXdNREF3TURJd05EaGhOalpsWVRkaE5URTVObUk0TW1Fell6VTRZbVJrWkdJMk1UQTBNV0l3WXpZNE16Z3dNREJrWkdZMFpUaGlaalEzTm1FNFpXSXpOV0UxWlRjeU1XUm1NbVpoTXpsak5qRmpORFl5WVRGa056VmxZMk5sTWpBME1tUmpOemcxSWl3aVpHVmpjbmx3ZEdsdmJrTnZibVJwZEdsdmJuTWlPaUpYTTNOcFdUSTVkV1JJU21oWk0xSkNXa2RTZVZwWVRucEphbTlwU1dsM2FXTXpVbWhpYlZKb1kyMVNSR0l5TlRCamJVWnFaRVpTTldOSFZXbFBhVWxwVEVOS2FtRkhSbkJpYVVrMlNXMVdNR0ZIVm5sYVdGWjBTV2wzYVdKWFZqQmhSemxyU1dwdmFVbHBkMmxqUjBaNVdWY3hiR1JIVm5samVVazJWM2xKTm1SWVRteGphMFpyV2toS2JHTXpUV2xZVTNkcFkyMVdNR1JZU25WV2JVWnpaRmRXVlZwWVRqQkphbkEzU1cxT2RtSllRbWhqYlVZd1lqTkphVTlwU1RsSmFYZHBaRzFHYzJSWFZXbFBhVWwzWlVSTmVFMXRWa0pQUkZWNVRucEpNbEpVVGtKUFYxa3lUWHBPUWsxRVRUTk9NazEzV2xkRk5FOUVTWGRQUkZwclRtcFpNazVxV1dsbVdEQnpaWGxLZG1OSFZubFpXRkoyWTJsSk5rbHRSblZhUTBvNVRFaHphVmt5T1hWa1NFcG9XVE5TUWxwSFVubGFXRTU2U1dwdmFVbHBkMmxqTTFKb1ltMVNhR050VWtSaU1qVXdZMjFHYW1SR1VqVmpSMVZwVDJsS1ZGTldaRVpKYVhkcFdUSm9hR0ZYTkdsUGFVcHNaRWRvYkdOdFZqRmlVMGx6U1cweGJHUkhhSFphUTBrMlNXbEpjMGx1UW1oamJVWjBXbGhTYkdOdVRXbFBiSE5wVDI1S2JHTXlPVEZqYlU1c1kzbEtaRXhEU25sYVdGSXhZMjAxVjFsWGVERmFWbEpzWXpOUmFVOXVjMmxaTWpsMFkwZEdlVmxZVW5aamFVazJTVzFPZG1KdVVtaGhWelY2U1dsM2FXUnRSbk5rVjFWcFQybEthbHBZU21oaVYyeHFUMms0ZGt0cU9YUmlNbEpzWWtReGNtRnVjSE5PYldneVdtNUthV1I2V21wT1dFWnJaVzVrY0U5WFZucGxTRm93VFZoWk1XSllVakJPTWpsclRqSm9hVTFxYXpCT2VsbDVUa2N4ZFU1SVZYZGpiVEY0VFZoS2IwOVhSblZoYlU1MVpVaG5hV1pZTUhObGVVcDJZMGRXZVZsWVVuWmphVWsyU1cxR2RWcERTamxNU0hOcFdUSTVkV1JJU21oWk0xSkNXa2RTZVZwWVRucEphbTlwU1dsM2FXTXpVbWhpYlZKb1kyMVNSR0l5TlRCamJVWnFaRVpTTldOSFZXbFBhVXBVVTFaa1JrbHBkMmxaTW1ob1lWYzBhVTlwU214a1IyaHNZMjFXTVdKVFNYTkpiVEZzWkVkb2RscERTVFpKYVVselNXNUNhR050Um5SYVdGSnNZMjVOYVU5c2MybFBia3BzWXpJNU1XTnRUbXhqZVVwa1RFTktlVnBZVWpGamJUVlhXVmQ0TVZwV1VteGpNMUZwVDI1emFWa3lPWFJqUjBaNVdWaFNkbU5wU1RaSmJVNTJZbTVTYUdGWE5YcEphWGRwWkcxR2MyUlhWV2xQYVVwcVdsaEthR0pYYkdwUGFUaDJTMm81ZEdJeVVteGlSREZ5WVc1d2MwNXRhREphYmtwcFpIcGFhazV0Um10T00yeHJZbXBDYjJGVVVqSmtSMFowWlVSS01rNXFTWGRoUjFKdVpGUmFjMkpJUlRCUFYyZDVUMGhLYlZwRVdtcGpla0Y1V25wT2FtSlhORFZsYlVWcFpsZ3djMlY1U25aalIxWjVXVmhTZG1OcFNUWkpiVVoxV2tOS09VeEljMmxaTWpsMVpFaEthRmt6VWtKYVIxSjVXbGhPZWtscWIybEphWGRwWXpOU2FHSnRVbWhqYlZKRVlqSTFNR050Um1wa1JsSTFZMGRWYVU5cFNsUlRWbVJHU1dsM2FWa3lhR2hoVnpScFQybEtiR1JIYUd4amJWWXhZbE5KYzBsdE1XeGtSMmgyV2tOSk5rbHBTWE5KYmtKb1kyMUdkRnBZVW14amJrMXBUMnh6YVU5dVNteGpNamt4WTIxT2JHTjVTbVJNUTBwNVdsaFNNV050TlZkWlYzZ3hXbFpTYkdNelVXbFBibk5wV1RJNWRHTkhSbmxaV0ZKMlkybEpOa2x0VG5aaWJsSm9ZVmMxZWtscGQybGtiVVp6WkZkVmFVOXBTbXBhV0Vwb1lsZHNhazlwT0haTGFqbDBZakpTYkdKRU1YSmhibkJ6VG0xb01scHVTbWxrZWxwcVRucFplbVJYU210aFJ6a3paVzFHZGsxSE1EQmxXRUUwVGtkT05HVnRTbTFpYlhodlRrZG9hMkZVVm1oaVNFWjJUa2hzZVZwWFNuUlpla0o0WTBkd2EyRlVWV2xtV0RGa0lpd2laVzVqY25sd2RHVmtJam9pUVVSWWNXYzNXbGRLVkVwNVNYUllaakZpTTJ0T01YcFliVFJCTVZoQlJVdzRXRXRUZVRWdFdYRnRaMlJ1UzI5RVprTmlWM2RoVFhkNWQwWlRlVWRtTmxaVlZGaHJlWEZ6ZW1sNFVrSkhMV3d6T0d4cFRWVmZOREp3WldJME0yWjBZa2RTYjJOblRXRk1NblZMTWpadFVYSXpNVlJDY0dSWlFWbDRSbGxWVnpVaWZRaWNyZWF0ZWRBdHgYMjAyMy0wNi0xMVQxMToyNDoyMy4wMjlaaXVwZGF0ZWRBdHgYMjAyMy0wNi0xMVQxMToyNDoyMy4wMjlaamFwcFZlcnNpb25lMC4yLjBqZm9sZGVyVHlwZQFwY29udGVudEZvbGRlcklkc4FkdGVtcGZoZWFkZXKjZW1vZGVsWCjOAQIBhQESIBzfRtp3sMC9cDyAKBAj8zq/N2Xd5IJPzEEoVKyvyJM1ZnVuaXF1ZUwrbXEamTAna5uAMqJrY29udHJvbGxlcnOBeDtkaWQ6cGtoOmVpcDE1NToxOjB4MzEyZUE4NTI3MjZFM0E5ZjYzM0EwMzc3YzBlYTg4MjA4NmQ2NjY2Ng=="},
			wantCommit: &ceramic.GenesisCommitPayload{
				Header: ceramic.CommitHeader{
					Model: ceramic.MustParseStreamID("kjzl6hvfrbw6c5qdzwi9esxvt1v5mtt7od7hb2947624mn4u0rmq1rh9anjcnxx"),
					Controllers: []string{
						"did:pkh:eip155:1:0x312eA852726E3A9f633A0377c0ea882086d66666",
					},
					Unique: []byte{43, 109, 113, 26, 153, 48, 39, 107, 155, 128, 50, 162},
				},
				Data: []byte(`{"options":"eyJmb2xkZXJOYW1lIjoiVW50aXRsZWQiLCJmb2xkZXJEZXNjcmlwdGlvbiI6IiIsImVuY3J5cHRlZFN5bW1ldHJpY0tleSI6ImM3ZTZmNDk4ODUyZWYyZWM4MTQwYzZhYjcwZDdjYWU4YmQwZjM5NGQ4NzhjMDM4OTNjZGRiMWU5YjlhMzA1MTJhZGEyZDUyNzczZTJmNWM5YzBkMTU0MjNiMjNhMmQwZWRkZGU5YTczYjM0Y2JlMmM3MDRkMjI2MTY0MzFhNDA3ZDUzMzIxZDVhNDljZGRiNTc2Y2I2OGNlYTcwMmE2NmZmNmZmNGQ4OGNkOWZlZjM0Zjc1MTNhNzZkZjhmNzA4ZmZkMDBiZWI3YjAwNmI4ZDY0OTUwNzlkZThiZTAzZGY3MDU1NjRmNDBjM2JjYWE1MWU2MGE3MDBiOGMzYmZlYzEwMDAwMDAwMDAwMDAwMDIwNDhhNjZlYTdhNTE5NmI4MmEzYzU4YmRkZGI2MTA0MWIwYzY4MzgwMDBkZGY0ZThiZjQ3NmE4ZWIzNWE1ZTcyMWRmMmZhMzljNjFjNDYyYTFkNzVlY2NlMjA0MmRjNzg1IiwiZGVjcnlwdGlvbkNvbmRpdGlvbnMiOiJXM3NpWTI5dWRISmhZM1JCWkdSeVpYTnpJam9pSWl3aWMzUmhibVJoY21SRGIyNTBjbUZqZEZSNWNHVWlPaUlpTENKamFHRnBiaUk2SW1WMGFHVnlaWFZ0SWl3aWJXVjBhRzlrSWpvaUlpd2ljR0Z5WVcxbGRHVnljeUk2V3lJNmRYTmxja0ZrWkhKbGMzTWlYU3dpY21WMGRYSnVWbUZzZFdWVVpYTjBJanA3SW1OdmJYQmhjbUYwYjNJaU9pSTlJaXdpZG1Gc2RXVWlPaUl3ZURNeE1tVkJPRFV5TnpJMlJUTkJPV1kyTXpOQk1ETTNOMk13WldFNE9ESXdPRFprTmpZMk5qWWlmWDBzZXlKdmNHVnlZWFJ2Y2lJNkltRnVaQ0o5TEhzaVkyOXVkSEpoWTNSQlpHUnlaWE56SWpvaUlpd2ljM1JoYm1SaGNtUkRiMjUwY21GamRGUjVjR1VpT2lKVFNWZEZJaXdpWTJoaGFXNGlPaUpsZEdobGNtVjFiU0lzSW0xbGRHaHZaQ0k2SWlJc0luQmhjbUZ0WlhSbGNuTWlPbHNpT25KbGMyOTFjbU5sY3lKZExDSnlaWFIxY201V1lXeDFaVlJsYzNRaU9uc2lZMjl0Y0dGeVlYUnZjaUk2SW1OdmJuUmhhVzV6SWl3aWRtRnNkV1VpT2lKalpYSmhiV2xqT2k4dktqOXRiMlJsYkQxcmFucHNObWgyWm5KaWR6WmpOWEZrZW5kcE9XVnplSFowTVhZMWJYUjBOMjlrTjJoaU1qazBOell5TkcxdU5IVXdjbTF4TVhKb09XRnVhbU51ZUhnaWZYMHNleUp2Y0dWeVlYUnZjaUk2SW1GdVpDSjlMSHNpWTI5dWRISmhZM1JCWkdSeVpYTnpJam9pSWl3aWMzUmhibVJoY21SRGIyNTBjbUZqZEZSNWNHVWlPaUpUU1ZkRklpd2lZMmhoYVc0aU9pSmxkR2hsY21WMWJTSXNJbTFsZEdodlpDSTZJaUlzSW5CaGNtRnRaWFJsY25NaU9sc2lPbkpsYzI5MWNtTmxjeUpkTENKeVpYUjFjbTVXWVd4MVpWUmxjM1FpT25zaVkyOXRjR0Z5WVhSdmNpSTZJbU52Ym5SaGFXNXpJaXdpZG1Gc2RXVWlPaUpqWlhKaGJXbGpPaTh2S2o5dGIyUmxiRDFyYW5wc05taDJabkppZHpaak5tRmtOM2xrYmpCb2FUUjJkR0Z0ZURKMk5qSXdhR1JuZFRac2JIRTBPV2d5T0hKbVpEWmpjekF5WnpOamJXNDVlbUVpZlgwc2V5SnZjR1Z5WVhSdmNpSTZJbUZ1WkNKOUxIc2lZMjl1ZEhKaFkzUkJaR1J5WlhOeklqb2lJaXdpYzNSaGJtUmhjbVJEYjI1MGNtRmpkRlI1Y0dVaU9pSlRTVmRGSWl3aVkyaGhhVzRpT2lKbGRHaGxjbVYxYlNJc0ltMWxkR2h2WkNJNklpSXNJbkJoY21GdFpYUmxjbk1pT2xzaU9uSmxjMjkxY21ObGN5SmRMQ0p5WlhSMWNtNVdZV3gxWlZSbGMzUWlPbnNpWTI5dGNHRnlZWFJ2Y2lJNkltTnZiblJoYVc1eklpd2lkbUZzZFdVaU9pSmpaWEpoYldsak9pOHZLajl0YjJSbGJEMXJhbnBzTm1oMlpuSmlkelpqTnpZemRXSmthRzkzZW1Gdk1HMDBlWEE0TkdONGVtSm1ibXhvTkdoa2FUVmhiSEZ2TkhseVpXSnRZekJ4Y0dwa2FUVWlmWDFkIiwiZW5jcnlwdGVkIjoiQURYcWc3WldKVEp5SXRYZjFiM2tOMXpYbTRBMVhBRUw4WEtTeTVtWXFtZ2RuS29EZkNiV3dhTXd5d0ZTeUdmNlZVVFhreXFzeml4UkJHLWwzOGxpTVVfNDJwZWI0M2Z0YkdSb2NnTWFMMnVLMjZtUXIzMVRCcGRZQVl4RllVVzUifQ","createdAt":"2023-06-11T11:24:23.029Z","updatedAt":"2023-06-11T11:24:23.029Z","appVersion":"0.2.0","folderType":1,"contentFolderIds":["temp"]}`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commit := &ceramic.GenesisCommitPayload{}
			buf := bytes.NewBuffer(lo.Must(base64.StdEncoding.DecodeString(tt.args.base64Encoded)))
			nd, _ := ceramic.DecodeDagCborNodeDataFromReader(buf)
			if err := commit.DecodeFromNodeData(nd); (err != nil) != tt.wantErr {
				t.Errorf("GenesisCommit.DecodeFromBlockReader() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(commit, tt.wantCommit) {
				t.Errorf("ConvertFrom() = %v\nwant %v", commit, tt.wantCommit)
			}
		})
	}
}

func TestAnchorCommit_DecodeFromBlockReader(t *testing.T) {
	type args struct {
		base64Encoded string
	}
	tests := []struct {
		name       string
		wantCommit *ceramic.AnchorCommit
		args       args
		wantErr    bool
	}{
		{
			name: "common",
			wantCommit: &ceramic.AnchorCommit{
				ID:    cid.MustParse("bagcqcera73sgdmuyznkpycnrkskk222l7qu6menvrx2ldyenjxdmsdabru6q"),
				Prev:  cid.MustParse("bagcqcerafrbuvb252ortgwwdpequn6i3bn67qxiholbff2irmqgqp6bmtxuq"),
				Proof: cid.MustParse("bafyreidtdpcjnltl7enswtp4s4xbsweb5zndvzihiyczl3t6ppqvbcgjpu"),
				Path:  "0/0/0/1/0/0/0/0/1",
			},
			args: args{
				base64Encoded: "pGJpZNgqWCYAAYUBEiD+5GGymMtU/AmxVJSta0v8KeYRtY30seCNTcbJDAGNPWRwYXRocTAvMC8wLzEvMC8wLzAvMC8xZHByZXbYKlgmAAGFARIgLENKh13TozNaw3khRvkbC334XQdywlLpEWQNB/gsnellcHJvb2bYKlglAAFxEiBzG8SWrmv5GytN/JcuGViB7lo65QdGBZXufnvhUIjJfQ==",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commit := &ceramic.AnchorCommit{}
			buf := bytes.NewBuffer(lo.Must(base64.StdEncoding.DecodeString(tt.args.base64Encoded)))
			nd, _ := ceramic.DecodeDagCborNodeDataFromReader(buf)
			if err := commit.DecodeFromNodeData(nd); (err != nil) != tt.wantErr {
				t.Errorf("AnchorCommit.DecodeFromBlockReader() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(commit, tt.wantCommit) {
				t.Errorf("ConvertFrom() = %v\nwant %v", commit, tt.wantCommit)
			}
		})
	}
}

func TestSignedCommit_DecodeFromBlockReader(t *testing.T) {
	type args struct {
		base64Encoded string
	}
	tests := []struct {
		name       string
		wantCommit *ceramic.SignedCommit
		args       args
		wantErr    bool
	}{
		{
			name: "common",
			args: args{
				base64Encoded: "omdwYXlsb2FkWCQBcRIgGtqDmUTOFzN1JfGKAS9bP1KGrprc6FzbjWjnGQv9yjlqc2lnbmF0dXJlc4GiaXByb3RlY3RlZFjMeyJhbGciOiJFZERTQSIsImNhcCI6ImlwZnM6Ly9iYWZ5cmVpZm5qdno3NDZ1azJrcWhpdHpqN2I1b3hhNmRjcDR3N2J1ZzJhbnhqbHNlYXhiM2R0bG9qeSIsImtpZCI6ImRpZDprZXk6ejZNa3UxOUdnY2gxTlliUmlOREZQVGFwcmQ0aUZKVlVncThTTnpNNnNSTEZESk45I3o2TWt1MTlHZ2NoMU5ZYlJpTkRGUFRhcHJkNGlGSlZVZ3E4U056TTZzUkxGREpOOSJ9aXNpZ25hdHVyZVhAoy7KX9UhdAIpYq2cvN8Kjv+Ka3ps6GzjYHt6afsu8qYYI7pS4k9/Ln013MBMnUTe7eM0/+YuAKnnOmsO8a5ECw==",
			},
			wantCommit: &ceramic.SignedCommit{
				Payload: cid.MustParse("bafyreia23kbzsrgoc4zxkjprrias6wz7kkdk5gw45bonxdli44mqx7okhe"),
				Signatures: []ceramic.Signature{
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
			commit := &ceramic.SignedCommit{}
			buf := bytes.NewBuffer(lo.Must(base64.StdEncoding.DecodeString(tt.args.base64Encoded)))
			nd, _ := ceramic.DecodeDagJWSNodeDataFromReader(buf)
			if err := commit.DecodeFromNodeData(nd); (err != nil) != tt.wantErr {
				t.Errorf("SignedCommit.DecodeFromBlockReader() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(commit, tt.wantCommit) {
				t.Errorf("ConvertFrom() = %v\nwant %v", commit, tt.wantCommit)
			}
		})
	}
}
