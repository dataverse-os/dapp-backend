package routers

import (
	"github.com/dataverse-os/dapp-backend/internal/dapp"
	"github.com/gin-gonic/gin"
)

type GetStreamFileBindingForm struct {
	ModelID          string  `form:"modelId" binding:"required"`
	IndexFileModelID string  `form:"indexFileModelId" binding:"required"`
	PKH              *string `form:"pkh"`
}

func GetStreamFileBinding(ctx *gin.Context) {
	var form GetStreamFileBindingForm
	if err := ctx.BindQuery(&form); err != nil {
		ctx.AbortWithStatusJSON(400, &gin.H{
			"msg": err.Error(),
		})
		return
	}

	result, err := dapp.ListStreamBindingFiles(ctx, form.ModelID, form.IndexFileModelID, form.PKH)
	if err != nil {
		ctx.AbortWithStatusJSON(400, &gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(200, result)
}
