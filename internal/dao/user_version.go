package dao

import (
	"context"
	"dapp-backend/internal"
	"dapp-backend/model"
	"dapp-backend/verify"
)

func UpdateFsVersion(ctx context.Context, origin, signed string) (res bool, err error) {
	var userVersion model.UserVsion
	if userVersion.Did, err = verify.ExportPublicKey(origin, signed); err != nil {
		return
	}
	userVersion.FsVersion = origin

	// query := internal.DB.Model(&model.UserVsion{}).Where(&model.UserVsion{Did: userVersion.Did}).Updates(userVersion)
	// if err = query.Error; err != nil {
	// return
	// }
	return
}

func GetUserVersionList(ctx context.Context, offset, limit *int) (list []model.UserVsion, err error) {
	query := internal.DB
	if offset != nil {
		query = query.Offset(*offset)
	}
	if limit != nil && *limit < 100 {
		query = query.Limit(*limit)
	} else {
		// query = query.Limit(20)
		query = query.Limit(3000)
	}
	err = query.Find(&list).Error
	return
}

func CreateUserVersion(ctx context.Context, origin, signed string) (list []model.UserVsion, err error) {
	var userVersion model.UserVsion
	// Check signed message
	if userVersion.Did, err = verify.ExportPublicKey(origin, signed); err != nil {
		return
	}
	err = internal.DB.Create(&userVersion).Error
	return
}
