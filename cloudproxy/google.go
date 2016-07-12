package cloudproxy

import (
	"io/ioutil"
	"net/http"

	"github.com/elazarl/goproxy"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func injectGoogleToken(ts oauth2.TokenSource) goproxy.ReqHandler {
	return goproxy.FuncReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		token, err := ts.Token()
		if err != nil {
			logger.Warn("Error getting Google oauth2 token ", err)
		}
		logger.Printf("Adding Google oauth2 headers (token %s expires at %s)", token.TokenType,
			token.Expiry)
		token.SetAuthHeader(req)
		return req, nil
	})
}

func (p CloudProxy) ConfigureGoogle(credentialFile string) {
	logger.Info("Reading Google credentials from ", credentialFile)
	data, err := ioutil.ReadFile(credentialFile)
	if err != nil {
		logger.Fatal("Error reading Google credentials: ", err)
	}
	conf, err := google.JWTConfigFromJSON(data,
		"https://www.googleapis.com/auth/bigquery",
		"https://www.googleapis.com/auth/devstorage.read_write")
	if err != nil {
		logger.Fatal("Error parsing JWT config from JSON: ", err)
	}
	//client := conf.Client(oauth2.NoContext) // TODO: can we get a token without making a client?
	tokenSource := conf.TokenSource(oauth2.NoContext)

	p.Proxy.OnRequest(goproxy.DstHostIs("www.googleapis.com:443")).Do(injectGoogleToken(tokenSource))
}
