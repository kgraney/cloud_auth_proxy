package cloudproxy

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var logger = log.WithFields(log.Fields{
	"proxy": "cloudproxy",
})

func CloudProxyMain(c *cli.Context) {
	port := uint(c.GlobalInt("port"))

	logger.Infof("Cloud proxy will listen on %d", port)

	proxy := NewCloudProxy(port)
	proxy.Configure()

	krbCredential := c.String("krb-credential")
	if krbCredential != "" {
		logger.Infof("Using Kerberos credential %s", krbCredential)
		proxy.ConfigureKerberos(krbCredential)
	}
	proxy.ListenAndServe(port)
}
