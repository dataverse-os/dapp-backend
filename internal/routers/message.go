package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

type ResponseNonce[T any] struct {
	Message string
	Data    T
	Nonce   string
}

func ResponseError(ctx *gin.Context, err error, code int) {
	resp := ResponseNonce[string]{
		Message: "Success",
		Nonce:   ctx.GetHeader("dataverse-nonce"),
		Data:    err.Error(),
	}
	ctx.Render(code, yamlRender{render.YAML{Data: resp}})
	ctx.Abort()
}
