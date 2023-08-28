package dapptable

import (
	"context"
	"testing"
)

func TestGetDappByModelID(t *testing.T) {
	type args struct {
		ctx     context.Context
		modelId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
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
			gotDapp, err := GetDappByModelID(tt.args.ctx, tt.args.modelId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDappByModelID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			isRightDapp := false
			for _, app := range gotDapp.Models {
				for _, model := range app.Streams {
					if model.ModelId == tt.args.modelId {
						isRightDapp = true
					}
				}
			}
			if !isRightDapp {
				t.Errorf("GetDappByModelID().Models = %v, want contain %v", gotDapp.Models, tt.args.modelId)
			}
		})
	}
}
