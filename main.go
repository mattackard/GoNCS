package main

import (
	"encoding/json"
	_ "expvar"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	_ "net/http/pprof"
	"net/url"
)

//Port server port
var Port string

//URL1 optional server location 1
var URL1 string

//URL2 optional server location 2
var URL2 string

//RequestCondition holds the request condition for proxy routing
type RequestCondition struct {
	RequestCondition string `json:"requestCondition"`
}

func init() {
	Port = ":6060"
	URL1 = "http://localhost:6061"
	URL2 = "http://localhost:6062"
}

func main() {

	http.HandleFunc("/", redirectRequest)

	fmt.Printf("Server running on localhost%s\n", Port)
	http.ListenAndServe(Port, nil)
}

func serveReverseProxy(target string, w http.ResponseWriter, r *http.Request) {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatal(err, "ln 47")
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	//set forwarding headers and values for forwarding request
	r.URL.Host = targetURL.Host
	r.URL.Scheme = targetURL.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = targetURL.Host

	proxy.ServeHTTP(w, r)
}

//checks the request's proxy condition and returns the url to forward to
func setProxyURL(proxyCondition string) string {
	if proxyCondition == "1" {
		return URL1
	}
	return URL2
}

//reads the request body and returns the proxy condition
func parseResponseBody(r *http.Request) RequestCondition {
	condition := RequestCondition{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err, "ln 74")
	}

	json.Unmarshal(body, &condition)
	return condition
}

//gathers the proxy condition and URL to forward the request
func redirectRequest(w http.ResponseWriter, r *http.Request) {
	condition := parseResponseBody(r)
	url := setProxyURL(condition.RequestCondition)
	println(url)
	serveReverseProxy(url, w, r)
}
