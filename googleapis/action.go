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

// Starts a listening reverse proxy to googleapis.
func ReverseProxyMain(c *cli.Context) {
	port := c.GlobalInt("port")

	logger.Infof("Reverse proxy to %s://%s will listen on %d", Scheme, RemoteHostname, port)

	http.HandleFunc("/discovery/", discoveryHandler)
	http.HandleFunc("/", proxyHandler(BuildReverseProxy(c.String("credentials"))))
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", port),
		c.GlobalString("certfile"),
		c.GlobalString("keyfile"),
		nil))
}
