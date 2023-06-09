package routers

import (
	"bytes"
	"crypto/ecdsa"
	"dapp-backend/verify"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
	key    *ecdsa.PrivateKey
)

func InitRouter() {
	var err error
	if key, err = crypto.HexToECDSA(os.Getenv("DID_PRIVATE_KEY")); err != nil {
		log.Fatalln(err)
	}

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
	router.Any("/api/*path", func(ctx *gin.Context) {
		u, _ := url.Parse(os.Getenv("CERAMIC_URL"))
		ctx.Request.URL.Scheme = u.Scheme
		ctx.Request.URL.Host = u.Host
		req := httputil.NewSingleHostReverseProxy(ctx.Request.URL)
		fmt.Println("send to: ", ctx.Request.URL)
		req.ServeHTTP(ctx.Writer, ctx.Request)
		ctx.Abort()
	})
	d := router.Group("/dataverse", checkWithNonce, CheckMiddleware())
	{
		d.POST("/validate", validate)
		d.POST("/dapp", createDapp)
		d.POST("/model", createModel)
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
		ctx.AbortWithStatusJSON(400, fmt.Errorf("invalid nonce"))
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
			ctx.AbortWithStatusJSON(400, err)
			return
		}
		if err = verify.CheckSign(data.Bytes(), ctx.GetHeader("dataverse-sig"), &key.PublicKey); err != nil {
			ctx.AbortWithStatusJSON(403, err)
			return
		}
		ctx.Request.Body = io.NopCloser(&data)
	}
}
