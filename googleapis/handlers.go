package googleapis

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func BuildReverseProxy() *httputil.ReverseProxy {
	// TODO: renew oauth2 tokens
	data, err := ioutil.ReadFile("/home/kevin/Google Sandbox-42115b17cf96.json")
	if err != nil {
		logger.Fatal(err)
	}
	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/bigquery")
	if err != nil {
		logger.Fatal(err)
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
