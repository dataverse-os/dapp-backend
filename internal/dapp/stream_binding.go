package dapp

import (
	"context"

	"github.com/dataverse-os/dapp-backend/ceramic"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
	"gorm.io/datatypes"
)

type StreamFileBinding struct {
	PKH           string        `json:"pkh"`
	ModelID       string        `json:"modelId"`
	StreamContent StreamContent `json:"streamContent"`
}

type StreamContent struct {
	Content  datatypes.JSON `json:"content,omitempty"`
	File     datatypes.JSON `json:"file,omitempty"`
	FileID   string         `json:"fileId,omitempty"`
	Verified bool           `json:"verified,omitempty"`
}

func ListStreamBindingFiles(ctx context.Context, modelID, indexFileModelID string, pkh *string) (
	result map[string]StreamFileBinding,
	err error,
) {
	query := &ceramic.Composite{}
	if pkh != nil {
		query.ControllerDID = *pkh
	}
	var streams []ceramic.Composite
	if err = ceramic.ComposeDB.WithContext(ctx).
		Table(modelID).Where(query).Find(&streams).Error; err != nil {
		return
	}
	var indexFiles []ceramic.Composite
	if err = ceramic.ComposeDB.WithContext(ctx).
		Table(indexFileModelID).Where(query).
		Where("stream_content->>'contentId' in (?)", lo.Map(streams, func(item ceramic.Composite, _ int) string {
			return item.StreamID
		})).
		Find(&indexFiles).Error; err != nil {
		return
	}
	indexFilesMap := lo.SliceToMap(indexFiles, func(item ceramic.Composite) (string, ceramic.Composite) {
		return gjson.GetBytes(item.StreamContent, "contentId").String(), item
	})
	result = make(map[string]StreamFileBinding)
	for _, v := range streams {
		result[v.StreamID] = StreamFileBinding{
			PKH:     v.ControllerDID,
			ModelID: modelID,
			StreamContent: StreamContent{
				Content: v.StreamContent,
				File:    indexFilesMap[v.StreamID].StreamContent,
				FileID:  indexFilesMap[v.StreamID].StreamID,
			},
		}
	}
	return
}
