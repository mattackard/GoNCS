package main

import (
	"encoding/base64"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
)

//PORT holds the port number that the reverse proxy will be deployed to
var port = os.Getenv("PROXYPORT")

//redirect holds a string containing the url for the proxy to redirect to
var redirect = "http://noteserver:5555"

//holds the username used to identify requests from the proxy
var httpUser = os.Getenv("USERNAME")

//holds the password used to identify requests from the proxy
var httpPass = os.Getenv("USERNAME")

func main() {
	http.HandleFunc("/", redirectHandler)
	log.Fatalln(http.ListenAndServe(port, nil))
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	myURL := parseURL(redirect)
	forwardHeaders(r, myURL)

	//makes the request to the actual server
	response, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Fatal(err, "line 30")
	}

	//adds all response headers from server to response object being sent back to client
	for key, values := range response.Header {
		for _, value := range values {
			w.Header().Set(key, value)
		}
	}

	//write status code and copy body from server response
	w.WriteHeader(response.StatusCode)
	io.Copy(w, response.Body)
}

//converts a string URL into a *url.URL struct
func parseURL(target string) *url.URL {
	parsed, err := url.Parse(target)
	if err != nil {
		log.Fatal(err, "line 49")
	}
	return parsed
}

func forwardHeaders(r *http.Request, url *url.URL) {
	r.Host = url.Host
	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.RequestURI = ""

	//gets request's remote address without port number and sets it in forwarding header
	split, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Fatal(err, "line 63")
	}
	r.Header.Set("X-Forwarded-For", split)

	//add user and pass for authenticating with the main server
	auth := base64.StdEncoding.EncodeToString([]byte(os.Getenv("PROXYAUTH")))
	r.Header.Set("Proxy-Authorization", "Basic "+auth)
}
