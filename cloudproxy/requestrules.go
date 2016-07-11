package cloudproxy

import (
	"fmt"
	"net/http"

	"github.com/elazarl/goproxy"
)

// Block any sites other than allowed public cloud APIs from being served through this proxy
type BlockOtherSites struct{}

func (s BlockOtherSites) Conditions() []goproxy.ReqCondition {
	return []goproxy.ReqCondition{
		goproxy.Not(goproxy.DstHostIs("www.googleapis.com:443")),
		goproxy.Not(goproxy.DstHostIs("accounts.google.com:443")),
		goproxy.Not(goproxy.DstHostIs("dl.google.com:443")),
	}
}

func (s BlockOtherSites) Handle(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	logger.Print("Blocking access to host ", req.URL.Host)
	return req, goproxy.NewResponse(req,
		goproxy.ContentTypeText,
		http.StatusForbidden,
		fmt.Sprintf("Access to %s is not available through this proxy.", req.URL.Host))
}

// Strip internal headers (authorization, etc.)
type StripInternalHeaders struct{}

func (s StripInternalHeaders) Conditions() []goproxy.ReqCondition {
	return []goproxy.ReqCondition{}
}

func (s StripInternalHeaders) Handle(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	req.Header.Del("Authorization")
	logger.Print(req.Header)
	return req, nil
}
