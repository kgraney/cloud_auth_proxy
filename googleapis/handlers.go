package googleapis

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func BuildReverseProxy(credentialFile string) *httputil.ReverseProxy {
	// TODO: renew oauth2 tokens
	logger.Info("Constructing reverse proxy")

	logger.Info("Reading Google credentials from ", credentialFile)
	data, err := ioutil.ReadFile(credentialFile)
	if err != nil {
		logger.Fatal("Error reading Google credentials: ", err)
	}
	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/bigquery")
	if err != nil {
		logger.Fatal("Error parsing JWT config from JSON: ", err)
	}
	client := conf.Client(oauth2.NoContext) // TODO: can we get a token without making a client?

	director := func(req *http.Request) {
		target := url.URL{
			Scheme: Scheme,
			Host:   RemoteHostname,
		}
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
		//req.URL.Path =
	}
	return &httputil.ReverseProxy{
		Director:  director,
		Transport: client.Transport,
	}
}

func proxyHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Processing response from %s", r.URL)
		// Serve the response from Google transparently back.
		p.ServeHTTP(w, r)
	}
}
