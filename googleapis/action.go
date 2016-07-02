package googleapis

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var logger = log.WithFields(log.Fields{
	"proxy": "googleapis",
})

const Scheme string = "https"
const RemoteHostname string = "www.googleapis.com"

func Main(c *cli.Context) {
	port := c.GlobalInt("port")

	logger.Infof("Reverse proxy to %s://%s will listen on %d", Scheme, RemoteHostname, port)

	http.HandleFunc("/discovery/", discoveryHandler)
	http.HandleFunc("/", proxyHandler(BuildReverseProxy()))
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", port),
		"/home/kmg/server.pem",
		"/home/kmg/server.key",
		nil))
}
