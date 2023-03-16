package routers

import (
	"dapp-backend/dapp"
	"dapp-backend/internal"
	"dapp-backend/verify"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

func createDapp(ctx *gin.Context) {
	name := ctx.GetString("DATAVERSE_NAME")
	if name == "" {
		ctx.AbortWithError(400, fmt.Errorf("empty input dapp name"))
		return
	}
	schemas, err := CreateDappModels(ctx, name)
	if err != nil {
		ctx.AbortWithError(400, err)
		return
	}
	result, _ := json.Marshal(schemas)
	responseWithSignedNonce(ctx, []byte(fmt.Sprintf("Status:Success/n%s", result)))
}

func createModel(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithError(400, err)
		return
	}
	var modelPath, streamID string
	if streamID, modelPath, err = dapp.GenerateCompositeJsonWithGraphql(ctx, body); err != nil {
		ctx.AbortWithError(400, err)
		return
	}
	if err = dapp.DeployCompositeJson(ctx, modelPath); err != nil {
		ctx.AbortWithError(400, err)
		return
	}
	responseWithSignedNonce(ctx, []byte(fmt.Sprintf("Status:Success/n%s", streamID)))
}

func responseWithSignedNonce(ctx *gin.Context, data []byte) {
	nonce := ctx.GetHeader("dataverse-nonce")
	if nonce == "" {
		ctx.AbortWithError(400, fmt.Errorf("invalid nonce"))
		return
	}
	result := fmt.Sprintf("%s\nNonce:%s", data, nonce)
	sig, err := verify.SignData([]byte(result), internal.PrivateKey)
	if err != nil {
		ctx.AbortWithError(500, err)
	}
	ctx.Header("dataverse-sig", string(sig))
	ctx.String(200, result)
}
