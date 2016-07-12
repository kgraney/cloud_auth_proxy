package cloudproxy

import (
	"fmt"
	"net/http"

	"github.com/elazarl/goproxy"
)

type CloudProxy struct {
	Proxy         *goproxy.ProxyHttpServer
	Port          uint
	RequestChain  []RequestChainEntry
	ResponseChain []ResponseChainEntry
}

func NewCloudProxy(port uint) *CloudProxy {
	return &CloudProxy{
		Proxy:         goproxy.NewProxyHttpServer(),
		Port:          port,
		RequestChain:  BuildRequestChain(),
		ResponseChain: BuildResponseChain(),
	}
}

func (p CloudProxy) Configure() {
	for _, entry := range p.RequestChain {
		p.Proxy.OnRequest(entry.Conditions()...).DoFunc(entry.Handle)
	}

	for _, entry := range p.ResponseChain {
		p.Proxy.OnResponse(entry.Conditions()...).DoFunc(entry.Handle)
	}

	p.Proxy.Verbose = true
}

func (p CloudProxy) ListenAndServe(port uint) {
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", p.Port), p.Proxy))
}
