package routers

import (
	"net/http"

	"github.com/dataverse-os/dapp-backend/ceramic"
	"github.com/dataverse-os/dapp-backend/verify"
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

func deployDapp(ctx *gin.Context) {
	var msg DeployMessage
	if err := yaml.NewDecoder(ctx.Request.Body).Decode(&msg); err != nil {
		ctx.AbortWithStatusJSON(400, err)
		return
	}
	resp := ResponseNonce[map[string]ModelResult]{
		Message: "Success",
		Nonce:   ctx.GetHeader("dataverse-nonce"),
		Data:    make(map[string]ModelResult),
	}
	for _, v := range msg.Models {
		_, streamID, err := ceramic.GenerateComposite(ctx, v.Schema, CeramicURL, CeramicAdminKey)
		if err != nil {
			resp.Message = err.Error()
			break
		}
		resp.Data[streamID] = ModelResult{
			StreamID: streamID,
			Schema:   v.Schema,
		}
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

	sig, err := verify.SignData(bytes, ceramicAdminKey)
	if err != nil {
		return err
	}
	w.Header()["dataverse-sig"] = []string{sig}

	_, err = w.Write(bytes)
	return err
}
