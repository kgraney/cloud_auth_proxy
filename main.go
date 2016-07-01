package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/codegangsta/cli"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

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
	}
	return &httputil.ReverseProxy{
		Director:  director,
		Transport: client.Transport,
	}
}

func proxyHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Processing response from %s", r.URL)
		// Serve the response from Google transparently back.
		p.ServeHTTP(w, r)
	}
}

func discoveryHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Handling discovery request %s", req.URL)
	//log.Printf("%#v", *req)

	// Send a request to Google.
	newreq, _ := http.NewRequest(req.Method, "https://www.googleapis.com"+req.RequestURI, req.Body)
	newreq.Close = true
	newreq.Body = req.Body
	client := http.Client{}
	resp, err := client.Do(newreq)
	if err != nil {
		log.Printf("discoveryHandler Do: %s", err)
	}

	// Modify the response body as needed.
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var s map[string]interface{}
	if err := json.Unmarshal(body, &s); err != nil {
		log.Printf("JSON unmarshall: %s", err)
	}

	// For actual API documentation we change the hostname.
	if b, _ := regexp.MatchString("/discovery/v1/apis/.*/.*", req.URL.Path); b {
		s["baseUrl"] = strings.Replace(s["baseUrl"].(string), "www.googleapis.com", "localhost:10000", -1)
		s["rootUrl"] = strings.Replace(s["rootUrl"].(string), "www.googleapis.com", "localhost:10000", -1)
		s["auth"] = nil // TODO: is this respected by anything?
	}

	// For the index of APIs we also change the hostname.
	if b, _ := regexp.MatchString("/discovery/v1/apis(/)?$", req.URL.Path); b {
		for _, x := range s["items"].([]interface{}) {
			m := x.(map[string]interface{})
			m["discoveryRestUrl"] = strings.Replace(m["discoveryRestUrl"].(string), "www.googleapis.com", "localhost:10000", -1)
		}
	}

	// Write the modified response back to the client.
	// TODO: include headers from Google?
	newbody, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		log.Printf("JSON marshall: %s", err)
	}
	w.Write(newbody)
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

		http.HandleFunc("/discovery/", discoveryHandler)
		http.HandleFunc("/", proxyHandler(BuildReverseProxy()))
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", port), "/home/kevin/go/server.pem", "/home/kevin/go/server.key", nil))
	}
	app.Run(os.Args)
}
