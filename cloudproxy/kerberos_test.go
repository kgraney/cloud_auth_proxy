package cloudproxy

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNegotiateHeaderConditions(t *testing.T) {
	req, err := http.NewRequest("GET", "example.com", nil)
	if err != nil {
		t.Error(err)
	}

	f := HasNoProxyAuthHeader()

	// No header present by default
	_, inMap := req.Header["Proxy-Authorization"]
	assert.False(t, inMap, "Test expected no Proxy-Authorization header")
	assert.True(t, f(req, nil), "With no header the condition should be true")

	// Add the header and some base64 value
	req.Header.Set("Proxy-Authorization", fmt.Sprintf("Negotiate %s",
		base64.StdEncoding.EncodeToString([]byte("foobar"))))
	_, inMap = req.Header["Proxy-Authorization"]
	assert.True(t, inMap, "Test expected Proxy-Authorization header to be present")
	assert.False(t, f(req, nil), "With a 'Negotiate <base64>' header the condition should be false.")
}
