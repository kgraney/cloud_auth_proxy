package cloudproxy

import (
	"net/http"

	"github.com/elazarl/goproxy"
)

const versionString string = "CloudAuthProxy/0.0.1"

// Put debug headers into the response
type ResponseDebugHeaders struct{}

func (s ResponseDebugHeaders) Conditions() []goproxy.RespCondition {
	return []goproxy.RespCondition{}
}

func (s ResponseDebugHeaders) Handle(res *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	if res != nil {
		res.Header.Set("X-Two-Sigma-Proxy", versionString)
	}
	return res

}
