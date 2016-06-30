package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/codegangsta/cli"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	//"github.com/apcera/gssapi/spnego"
	//"github.com/elazarl/goproxy"
)

type mutateTransport struct {
	transport http.RoundTripper
}

func (t *mutateTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := t.transport.RoundTrip(request)

	body, err := httputil.DumpResponse(response, true)
	if err != nil {
		return nil, err
	}
	log.Print(string(body))

	return response, err
}

func BuildReverseProxy() *httputil.ReverseProxy {
	// TODO: renew oauth2 tokens

	data, err := ioutil.ReadFile("/home/kevin/Google Sandbox-42115b17cf96.json")
	if err != nil {
		log.Fatal(err)
	}
	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/bigquery")
	if err != nil {
		log.Fatal(err)
	}
	client := conf.Client(oauth2.NoContext) // TODO: can we get a token without making a client?

	director := func(req *http.Request) {
		target := url.URL{
			Scheme: "https",
			Host:   "www.googleapis.com",
		}
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
		//req.URL.Path =

		var buf bytes.Buffer
		req.Write(&buf)
		log.Print(&buf)
	}
	return &httputil.ReverseProxy{
		Director:  director,
		Transport: &mutateTransport{client.Transport},
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "cloud_proxy"
	app.Usage = "A forward/reverse proxy for Public Cloud APIs"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port",
			Usage: "The port to listen on",
		},
	}

	app.Action = func(c *cli.Context) {
		port := c.Int("port")
		scheme := "https"
		host := "www.googleapis.com"

		log.Printf("Reverse proxy to %s://%s will listen on %d", scheme, host, port)

		/*
			proxy := httputil.NewSingleHostReverseProxy(&url.URL{
				Scheme: scheme,
				Host:   host,
			})
		*/
		proxy := BuildReverseProxy()
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", port), "/home/kevin/go/server.pem", "/home/kevin/go/server.key", proxy))
	}
	app.Run(os.Args)
}
