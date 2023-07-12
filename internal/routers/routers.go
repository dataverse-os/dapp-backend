package routers

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dataverse-os/dapp-backend/internal/dapp"
	"github.com/dataverse-os/dapp-backend/verify"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func InitRouter() {
	router = gin.Default()
	router.Use(
		cors.New(cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{
				http.MethodGet, http.MethodPost, http.MethodPut,
				http.MethodDelete, http.MethodOptions, http.MethodPatch,
				http.MethodHead,
			},
			AllowHeaders:     []string{"Origin"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				return true
			},
			MaxAge: 12 * time.Hour,
		}),
	)
	router.Any("/api/*path", CeramicProxy)
	d := router.Group("/dataverse", checkWithNonce, CheckMiddleware())
	{
		d.POST("/validate", validate)
		d.POST("/dapp", deployDapp)
	}
}

func Start() {
	if err := router.Run(":8080"); err != nil {
		log.Fatalln(err)
	}
}

func checkWithNonce(ctx *gin.Context) {
	nonce := ctx.GetHeader("dataverse-nonce")
	if nonce == "" {
		ResponseError(ctx, errors.New("invalid nonce"), 400)
		return
	}
}

func CheckMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			data bytes.Buffer
			err  error
		)
		if _, err = io.Copy(&data, ctx.Request.Body); err != nil {
			ResponseError(ctx, err, 400)
			return
		}
		if err = verify.CheckSign(data.Bytes(), ctx.GetHeader("dataverse-sig"), &dapp.CeramicSession.AdminKey.PublicKey); err != nil {
			ResponseError(ctx, err, 403)
			return
		}
		ctx.Request.Body = io.NopCloser(&data)
	}
}
