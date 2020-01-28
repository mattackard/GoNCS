package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
)

//PORT holds the port number that the reverse proxy will be deployed to
var PORT = os.Getenv("PROXYPORT")

//REDIRECT holds a string containing the url for the proxy to redirect to
var REDIRECT = os.Getenv("REDIRECTURL")

func main() {
	http.HandleFunc("/", redirectHandler)
	http.ListenAndServe(PORT, nil)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	myURL := parseURL(REDIRECT)
	forwardHeaders(r, myURL)

	//gets request's remote address without port number and sets it in forwarding header
	split, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Set("X-Forwarded-For", split)

	//makes the request to the actual server
	response, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}
	return parsed
}

func forwardHeaders(r *http.Request, url *url.URL) {
	r.Host = url.Host
	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.RequestURI = ""
}
