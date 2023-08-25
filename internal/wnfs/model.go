package wnfs

import (
	"context"

	dapptable "github.com/dataverse-os/dapp-backend/pkg/dapp-table"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/samber/lo"
)

type DappModel struct {
	ModelID     string    `gorm:"primary-key"`
	AppID       uuid.UUID `gorm:"type:uuid;not null"`
	Encryptable pq.StringArray
	ModelName   string
	Version     int
}

func LookupDappModel(ctx context.Context, modelId string) (data DappModel, err error) {
	if err = db.Where(&DappModel{ModelID: modelId}).First(&data).Error; err == nil {
		return
	}
	var dapp dapptable.Dapp
	if dapp, err = dapptable.GetDappByModelID(ctx, modelId); err != nil {
		return
	}
	models := make([]DappModel, 0)
	for _, v := range dapp.Models {
		for _, vv := range v.Streams {
			models = append(models, DappModel{
				ModelID:     modelId,
				AppID:       uuid.MustParse(dapp.Id),
				Encryptable: vv.Encryptable,
				ModelName:   v.ModelName,
				Version:     vv.Version,
			})
		}
	}

	// lookup models not in database
	modelsInDB := make([]DappModel, 0)
	if err = db.Where("model_id in (?)", lo.Map(models, func(item DappModel, _ int) string {
		return item.ModelID
	})).Find(&modelsInDB).Error; err != nil {
		return
	}
	// filter models already in database
	modelsShouldInsert := lo.Filter(models, func(item DappModel, index int) bool {
		return lo.ContainsBy(modelsInDB, func(itemInDB DappModel) bool {
			return itemInDB.ModelID != item.ModelID
		})
	})
	if err = db.Create(&modelsShouldInsert).Error; err != nil {
		return
	}

	data = lo.Filter(models, func(item DappModel, _ int) bool {
		return item.ModelID == modelId
	})[0]
	return
}
