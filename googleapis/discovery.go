package googleapis

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

// We need to do some modification of the service discovery responses from the remote
// host.

func discoveryHandler(w http.ResponseWriter, req *http.Request) {
	logger.Printf("Handling discovery request %s", req.URL)
	//logger.Printf("%#v", *req)

	// Send a request to Google.
	newreq, _ := http.NewRequest(req.Method, "https://www.googleapis.com"+req.RequestURI, req.Body)
	newreq.Close = true
	newreq.Body = req.Body
	client := http.Client{}
	resp, err := client.Do(newreq)
	if err != nil {
		logger.Printf("discoveryHandler Do: %s", err)
	}

	// Modify the response body as needed.
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var s map[string]interface{}
	if err := json.Unmarshal(body, &s); err != nil {
		logger.Printf("JSON unmarshall: %s", err)
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
		logger.Printf("JSON marshall: %s", err)
	}
	w.Write(newbody)
}
