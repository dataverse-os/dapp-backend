package dapptable

import (
	"context"
	"reflect"
	"testing"
)

func TestGetDappByModelID(t *testing.T) {
	type args struct {
		ctx     context.Context
		modelId string
	}
	tests := []struct {
		name     string
		args     args
		wantDapp Dapp
		wantErr  bool
	}{
		{
			name: "common",
			args: args{
				ctx:     context.Background(),
				modelId: "kjzl6hvfrbw6c5m98besslbjufnwxk9t1uzebyu1gevzr17tq65sbe3vv8oq53b",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDappBytes, err := GetDappByModelID(tt.args.ctx, tt.args.modelId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDappByModelID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDappBytes, tt.wantDapp) {
				t.Errorf("GetDappByModelID() = %v, want %v", gotDappBytes, tt.wantDapp)
			}
		})
	}
}
