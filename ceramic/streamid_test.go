package ceramic_test

import (
	"reflect"
	"testing"

	"github.com/dataverse-os/dapp-backend/ceramic"
	"github.com/ipfs/go-cid"
)

func TestStreamID_String(t *testing.T) {
	type fields struct {
		Type ceramic.StreamIdType
		Cid  cid.Cid
		Log  cid.Cid
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "stream_id",
			fields: fields{
				Type: ceramic.StreamIdTypeModelInstanceDocument,
				Cid:  cid.MustParse("bagcqceragiwmxdtelb45wjl6calr45bh4rpcipbmh3t3skrkjvv3foihfmrq"),
			},
			want: "kjzl6kcym7w8y64sx1g97dy5v3xmm49mnx0p9mofs44t6y0y7wp2ko8z7w9azhf",
		},
		{
			name: "commit_id",
			fields: fields{
				Type: ceramic.StreamIdTypeTile,
				Cid:  cid.MustParse("bagcqcerakszw2vsovxznyp5gfnpdj4cqm2xiv76yd24wkjewhhykovorwo6a"),
				Log:  cid.MustParse("bagjqcgzaday6dzalvmy5ady2m5a5legq5zrbsnlxfc2bfxej532ds7htpova"),
			},
			want: "k1dpgaqe3i64kjqcp801r3sn7ysi5i0k7nxvs7j351s7kewfzr3l7mdxnj7szwo4kr9mn2qki5nnj0cv836ythy1t1gya9s25cn1nexst3jxi5o3h6qprfyju",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := ceramic.StreamID{
				Type: tt.fields.Type,
				Cid:  tt.fields.Cid,
				Log:  tt.fields.Log,
			}
			if got := id.String(); got != tt.want {
				t.Errorf("StreamID.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseStreamID(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		wantId  ceramic.StreamID
		wantErr bool
	}{
		{
			name: "stream_id",
			args: args{
				str: "kjzl6kcym7w8y64sx1g97dy5v3xmm49mnx0p9mofs44t6y0y7wp2ko8z7w9azhf",
			},
			wantId: ceramic.StreamID{
				Type: ceramic.StreamIdTypeModelInstanceDocument,
				Cid:  cid.MustParse("bagcqceragiwmxdtelb45wjl6calr45bh4rpcipbmh3t3skrkjvv3foihfmrq"),
			},
			wantErr: false,
		},
		{
			name: "commit_id",
			args: args{
				str: "k1dpgaqe3i64kjqcp801r3sn7ysi5i0k7nxvs7j351s7kewfzr3l7mdxnj7szwo4kr9mn2qki5nnj0cv836ythy1t1gya9s25cn1nexst3jxi5o3h6qprfyju",
			},
			wantId: ceramic.StreamID{
				Type: ceramic.StreamIdTypeTile,
				Cid:  cid.MustParse("bagcqcerakszw2vsovxznyp5gfnpdj4cqm2xiv76yd24wkjewhhykovorwo6a"),
				Log:  cid.MustParse("bagjqcgzaday6dzalvmy5ady2m5a5legq5zrbsnlxfc2bfxej532ds7htpova"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotId, err := ceramic.ParseStreamID(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStreamID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotId, tt.wantId) {
				t.Errorf("ParseStreamID() = %v, want %v", gotId, tt.wantId)
			}
		})
	}
}
