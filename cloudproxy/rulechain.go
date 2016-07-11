package cloudproxy

import (
	"net/http"

	"github.com/elazarl/goproxy"
)

type RequestChainEntry interface {
	Conditions() []goproxy.ReqCondition
	Handle(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response)
}

func BuildRequestChain() []RequestChainEntry {
	return []RequestChainEntry{
		&BlockOtherSites{},
		&StripInternalHeaders{},
	}
}

type ResponseChainEntry interface {
	Conditions() []goproxy.RespCondition
	Handle(res *http.Response, ctx *goproxy.ProxyCtx) *http.Response
}

func BuildResponseChain() []ResponseChainEntry {
	return []ResponseChainEntry{
		&ResponseDebugHeaders{},
	}
}
