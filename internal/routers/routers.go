package routers

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dataverse-os/dapp-backend/ceramic"
	"github.com/dataverse-os/dapp-backend/verify"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	router          *gin.Engine
	CeramicAdminKey = os.Getenv("DID_PRIVATE_KEY")
	ceramicAdminKey *ecdsa.PrivateKey
	CeramicURL      = os.Getenv("CERAMIC_URL")
	ceramicURL      *url.URL
	checkSign       = os.Getenv("NO_CHECK_SIGN") == ""
)

func init() {
	var err error
	if ceramicURL, err = url.Parse(CeramicURL); err != nil {
		log.Fatalf("cannot parse env CERAMIC_URL '%s' as url", CeramicURL)
	}

	if ceramicAdminKey, err = crypto.HexToECDSA(CeramicAdminKey); err != nil {
		log.Fatalf("failed to parse ceramic admin key with error: %s", err)
	}

	if err = ceramic.Default.CheckAdminAccess(context.Background(), CeramicURL, CeramicAdminKey); err != nil {
		log.Fatalf("failed to parse ceramic url with error: %s", err)
	}
}

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
		if checkSign {
			if err = verify.CheckSign(data.Bytes(), ctx.GetHeader("dataverse-sig"), &ceramicAdminKey.PublicKey); err != nil {
				ResponseError(ctx, err, 403)
				return
			}
		}
		ctx.Request.Body = io.NopCloser(&data)
	}
}
