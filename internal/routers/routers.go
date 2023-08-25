package routers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dataverse-os/dapp-backend/internal/dapp"
	"github.com/dataverse-os/dapp-backend/verify"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func InitRouter() {
	router = gin.Default()
	router.Any("/api/*path", CeramicProxy)
	d := router.Group("/dataverse")
	{
		d.POST("/validate", checkWithNonce, validate)
		if !dapp.IsSandbox {
			d.POST("/dapp", checkWithNonce, CheckMiddleware(), deployDapp)
		} else {
			d.POST("/dapp", checkWithNonce, deployDapp)
		}
		d.GET("/model-version", GetModelVersion)
		d.POST("/model-version", HeaderChecker("dataverse-sig"), SignatureMiddleware, PostUpdateModelVersion)

		if os.Getenv("ENABLE_STREAM_BINDING_FILE_INDEX") == "TRUE" {
			d.GET("/stream-binding", GetStreamFileBinding)
		}
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

func HeaderChecker(headers ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, v := range headers {
			if ctx.GetHeader(v) == "" {
				ctx.AbortWithStatusJSON(400, &gin.H{
					"msg": fmt.Sprintf("should contain header: %s", v),
				})
				return
			}
		}
	}
}

func SignatureMiddleware(ctx *gin.Context) {
	var (
		data            bytes.Buffer
		signatureString = ctx.GetHeader("dataverse-sig")
		address         common.Address
		err             error
	)
	if signatureString == "" {
		return
	}
	if _, err = io.Copy(&data, ctx.Request.Body); err != nil {
		ResponseError(ctx, err, 400)
		return
	}
	if address, err = verify.ExportAddress(data.Bytes(), signatureString); err != nil {
		ResponseError(ctx, err, 400)
		return
	} else {
		ctx.Set("DATAVERSE_ADDRESS", address)
	}
	ctx.Request.Body = io.NopCloser(&data)
}
