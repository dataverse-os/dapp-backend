package routers

import (
	"bytes"
	"net/http/httputil"
	"strings"

	"github.com/dataverse-os/dapp-backend/internal/dapp"
	"github.com/dataverse-os/dapp-backend/internal/wnfs"
	"github.com/gin-gonic/gin"
)

func CeramicProxy(ctx *gin.Context) {
	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.Out.Header = r.In.Header
			r.Out.Host = dapp.CeramicSession.URL.Host
			r.Out.URL.Host = dapp.CeramicSession.URL.Host
			r.Out.URL.Scheme = dapp.CeramicSession.URL.Scheme
		},
	}
	res := bodyLogWriter{
		ResponseWriter: ctx.Writer,
		body:           new(bytes.Buffer),
	}
	proxy.ServeHTTP(res, ctx.Request)
	for k := range wnfs.StreamResponsePath {
		if strings.HasPrefix(ctx.Request.URL.Path, k) {
			wnfs.AppendToVerifyGroup(res.body)
		}
	}
	ctx.Abort()
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
