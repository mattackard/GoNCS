package main

import (
	"encoding/base64"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/mattackard/project-1/pkg/logutil"
)

//PORT holds the port number services that communicate with the proxy
var proxyPort = os.Getenv("PROXYPORT")
var logPort = os.Getenv("LOGPORT")
var logName = os.Getenv("LOGGERNAME")
var serverPort = os.Getenv("SERVERPORT")
var serverName = os.Getenv("SERVERNAME")
var loggerAddr = logName + ":" + logPort
var serverAddr = "http://" + serverName + ":" + serverPort

//holds the user and pass used to identify requests from the proxy
var proxyAuth = os.Getenv("PROXYAUTH")

var logFile *os.File

func main() {
	logFile = logutil.OpenLogFile("./logs/")
	defer logFile.Close()
	http.HandleFunc("/", redirectHandler)
	log.Fatalln(http.ListenAndServe(":"+proxyPort, nil))
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	logutil.LogServerRequest(w, r, loggerAddr, nil, "ReverseProxy")
	myURL := parseURL(serverAddr)
	forwardHeaders(r, myURL)

	//makes the request to the actual server
	response, err := http.DefaultClient.Do(r)
	if err != nil {
		logutil.SendLog(loggerAddr, true, []string{err.Error(), "line 35"}, nil, "ReverseProxy")
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
		logutil.SendLog(loggerAddr, true, []string{err.Error(), "line 54"}, nil, "ReverseProxy")
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
		logutil.SendLog(loggerAddr, true, []string{err.Error(), "line 68"}, nil, "ReverseProxy")
	}
	r.Header.Set("X-Forwarded-For", split)

	//add user and pass for authenticating with the main server
	auth := base64.StdEncoding.EncodeToString([]byte(proxyAuth))
	r.Header.Set("Proxy-Authorization", "Basic "+auth)
}
