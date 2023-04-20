package routers

import (
	"dapp-backend/ceramic"
	"dapp-backend/verify"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"gopkg.in/yaml.v3"
)

func validate(ctx *gin.Context) {
	resp := ResponseNonce[any]{
		Message: "Success",
		Nonce:   ctx.GetHeader("dataverse-nonce"),
	}
	ctx.Render(200, yamlRender{render.YAML{Data: resp}})
}

func createDapp(ctx *gin.Context) {
	var msg CreateMessage
	if err := yaml.NewDecoder(ctx.Request.Body).Decode(&msg); err != nil {
		ctx.AbortWithStatusJSON(400, err)
		return
	}
	resp := ResponseNonce[map[string]string]{
		Message: "Success",
		Nonce:   ctx.GetHeader("dataverse-nonce"),
		Data:    make(map[string]string),
	}
	for _, v := range msg.Models {
		_, streamID, err := ceramic.GenerateComposite(ctx, v)
		if err != nil {
			resp.Message = err.Error()
			break
		}
		resp.Data[streamID] = v
	}
	ctx.Render(200, yamlRender{render.YAML{Data: resp}})
}

func createModel(ctx *gin.Context) {
	var msg SetExternalModelsMessage
	if err := yaml.NewDecoder(ctx.Request.Body).Decode(&msg); err != nil {
		ctx.AbortWithStatusJSON(400, err)
		return
	}
	resp := ResponseNonce[map[string]string]{
		Message: "Success",
		Nonce:   ctx.GetHeader("dataverse-nonce"),
	}
	_, streamID, err := ceramic.GenerateComposite(ctx, msg.Schema)
	if err != nil {
		resp.Message = err.Error()
	} else {
		resp.Data[streamID] = msg.Schema
	}
	ctx.Render(200, yamlRender{render.YAML{Data: resp}})
}

type yamlRender struct {
	render.YAML
}

func (r yamlRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	bytes, err := yaml.Marshal(r.Data)
	if err != nil {
		return err
	}

	sig, err := verify.SignData(bytes, key)
	if err != nil {
		return err
	}
	w.Header()["dataverse-sig"] = []string{hexutil.Encode(sig)}

	_, err = w.Write(bytes)
	return err
}
