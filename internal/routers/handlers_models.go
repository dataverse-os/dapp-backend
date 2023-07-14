package routers

import (
	"github.com/dataverse-os/dapp-backend/internal/dapp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

type GetModelVersionForm struct {
	Address   string `form:"address" binding:"required"`
	ModelName string `form:"modelName" binding:"required"`
}

func GetModelVersion(ctx *gin.Context) {
	var form GetModelVersionForm
	if err := ctx.BindQuery(&form); err != nil {
		ctx.AbortWithStatusJSON(400, &gin.H{
			"msg": err.Error(),
		})
		return
	}

	version, err := dapp.LookupUserModelVersion(common.HexToAddress(form.Address), form.ModelName)
	if err != nil {
		ctx.AbortWithStatusJSON(400, &gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(200, &gin.H{
		"version": version,
	})
}

type PostUpdateModelVersionForm struct {
	ModelName string `form:"modelName" binding:"required"`
	Version   uint64 `form:"version" binding:"required"`
}

func PostUpdateModelVersion(ctx *gin.Context) {
	var (
		form    PostUpdateModelVersionForm
		address common.Address
	)
	if err := ctx.BindJSON(&form); err != nil {
		ctx.AbortWithStatusJSON(400, &gin.H{
			"msg": err.Error(),
		})
		return
	}

	if val, ok := ctx.Get("DATAVERSE_ADDRESS"); ok && val != nil {
		if address, ok = val.(common.Address); !ok {
			ctx.AbortWithStatusJSON(400, &gin.H{
				"msg": "filed to parse into common.Address",
			})
			return
		}
	}
	if err := dapp.UpdateUserModelVersion(address, form.ModelName, form.Version); err != nil {
		ctx.AbortWithStatusJSON(400, &gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(200, &gin.H{
		"address":   address.Hex(),
		"modelName": form.ModelName,
		"version":   form.Version,
	})
}
