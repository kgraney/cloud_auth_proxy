package cloudproxy

import (
	"net"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/apcera/gssapi"
	"github.com/apcera/gssapi/spnego"
	"github.com/elazarl/goproxy"
)

var ProxyAuthConnectNegotiate = &goproxy.ConnectAction{
	Action:    goproxy.ConnectProxyAuthHijack,
	TLSConfig: goproxy.TLSConfigFromCA(&goproxy.GoproxyCa),
	Hijack: func(req *http.Request, client net.Conn, ctx *goproxy.ProxyCtx) {
		client.Write([]byte("Proxy-Authenticate: Negotiate\r\n\r\n"))
	},
}

func LoadKrb5Lib() (*gssapi.Lib, error) {
	options := gssapi.Options{
		LibPath:    "",
		Krb5Config: "/etc/krb5.conf",
		Krb5Ktname: os.Getenv("KRB5_KTNAME"),
	}

	lib, err := gssapi.Load(&options)
	return lib, err
}

func KrbProxyAuth(krbServer spnego.KerberizedServer, cred *gssapi.CredId) goproxy.ReqHandler {
	return goproxy.FuncReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		outHeader := http.Header{}
		principal, status, err := krbServer.Negotiate(cred, req.Header, outHeader)
		if err != nil {
			logger.Print("Proxy negotiate error: ", err)
		}

		authed_logger := logger.WithFields(log.Fields{
			"krb_principal": principal,
		})

		if status != 200 {
			resp := goproxy.NewResponse(req, goproxy.ContentTypeText, status,
				"Proxy Authentication Error")
			resp.Header = outHeader
			return req, resp
		}

		// Let the user through
		authed_logger.Print("Allowing access to URL ", ctx.Req.URL.String())
		return req, nil
	})
}

func KrbProxyAuthConnect(krbServer spnego.KerberizedServer, cred *gssapi.CredId) goproxy.HttpsHandler {
	return goproxy.FuncHttpsHandler(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		outHeader := http.Header{}
		principal, status, err := krbServer.Negotiate(cred, ctx.Req.Header, outHeader)
		if err != nil {
			logger.Print("Proxy negotiate error: ", err)
		}

		authed_logger := logger.WithFields(log.Fields{
			"krb_principal": principal,
		})

		if status != 200 {
			return ProxyAuthConnectNegotiate, host
		}

		// Let the user through
		authed_logger.Print("Allowing access to URL ", ctx.Req.URL.String())
		return goproxy.MitmConnect, host
	})
}

func IsNotHttpsRequest() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		return req.URL.Scheme != "https"
	}
}

func HasNoProxyAuthHeader() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		_, ok := req.Header["Proxy-Authorization"]
		return !ok
	}
}

func HasProxyNegotiateChallenge() goproxy.RespConditionFunc {
	return func(resp *http.Response, ctx *goproxy.ProxyCtx) bool {
		return (resp.StatusCode == 407) && (len(resp.Header["Proxy-Authenticate"]) > 0)
	}
}

func (p CloudProxy) ConfigureKerberos(krbCredential string) {
	lib, err := LoadKrb5Lib()
	if err != nil {
		logger.Print("Error loading Krb5 lib: ", err)
	}

	krbServer := spnego.KerberizedServer{
		Lib: lib,
		UseProxyAuthentication: true,
	}
	cred, err := krbServer.AcquireCred(krbCredential)
	if err != nil {
		logger.Print("Error creating KerberizedServer: ", err)
	}

	// Krb5 authentication is required on all CONNECT requests and all non-HTTPS requests.  HTTPS
	// requests are authenticated on the CONNECT.
	p.Proxy.OnRequest().HandleConnect(KrbProxyAuthConnect(krbServer, cred))
	p.Proxy.OnRequest(IsNotHttpsRequest()).Do(KrbProxyAuth(krbServer, cred))
}
