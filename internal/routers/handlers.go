package routers

import (
	"net/http"

	"github.com/dataverse-os/dapp-backend/internal/dapp"
	"github.com/dataverse-os/dapp-backend/verify"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/google/uuid"
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
	var (
		msg dapp.DeployMessage
		id  uuid.UUID
		err error
	)
	defer func() {
		if err != nil {
			ResponseError(ctx, err, 400)
		}
	}()
	if err = yaml.NewDecoder(ctx.Request.Body).Decode(&msg); err != nil {
		return
	}
	if id, err = uuid.Parse(ctx.GetHeader("dataverse-dapp-id")); err != nil {
		return
	}
	resp := ResponseNonce[[]dapp.ModelResult]{
		Message: "Success",
		Nonce:   ctx.GetHeader("dataverse-nonce"),
	}
	if resp.Data, err = dapp.DeployStreamModels(ctx, id, msg.Models, CeramicURL, CeramicAdminKey); err != nil {
		return
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
