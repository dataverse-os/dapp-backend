package routers

import (
	"dapp-backend/internal"
	sign "dapp-backend/verify"
	"log"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func InitRouter() {
	router = gin.Default()
	router.Any("/v0", func(ctx *gin.Context) {
		req := httputil.NewSingleHostReverseProxy(ctx.Request.URL)
		req.ServeHTTP(ctx.Writer, ctx.Request)
		ctx.Abort()
	})
	d := router.Group("/dataverse",
		sign.CheckMiddleware(internal.PrivateKey),
	)
	{
		d.POST("/dapp", sign.ParseBodyMiddleware, createDapp)
		d.POST("/model", createModel)
	}
}

func Start() {
	if err := router.Run(":8080"); err != nil {
		log.Fatalln(err)
	}
}
