package main

import (
	"encoding/base64"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/mattackard/project-1/pkg/dnsutil"
	"github.com/mattackard/project-1/pkg/logutil"
	"github.com/mattackard/project-1/pkg/perfutil"
)

//PORT holds the port number services that communicate with the proxy
var proxyPort = os.Getenv("PROXYPORT")

//holds the user and pass used to identify requests from the proxy
var proxyAuth = os.Getenv("PROXYAUTH")

var dnsPort = os.Getenv("DNSPORT")
var dnsName = os.Getenv("DNSNAME")
var dnsAddr = dnsName + ":" + dnsPort

var serverAddr string
var loggerAddr string
var logFile *os.File

func init() {
	//requests the address of services from the dns server
	serverAddr = dnsutil.GetServiceIP(dnsAddr, "noteserver")
	loggerAddr = dnsutil.GetServiceIP(dnsAddr, "logger")
}

func main() {
	logFile = logutil.OpenLogFile("./logs")
	defer logFile.Close()

	//send all requests through the redirect handler
	http.HandleFunc("/", redirectHandler)

	//handler for sending performance stats
	http.HandleFunc("/getStats", perfutil.SendStatsHTTP)

	//send messages to log file and terminal to record startup
	proxyIP := dnsutil.Ping(dnsAddr, "reverseproxy")
	logutil.SendLog(loggerAddr, false, []string{"Reverse Proxy Server started at " + dnsutil.TrimPort(proxyIP) + ":" + proxyPort}, logFile, "ReverseProxy")

	//start server
	err := http.ListenAndServe(":"+proxyPort, nil)
	if err != nil {
		logutil.SendLog("logger:6060", false, []string{"Reverse proxy at " + dnsutil.TrimPort(proxyIP) + " shutting down"}, logFile, "ReverseProxy")
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	logutil.LogServerRequest(w, r, loggerAddr, logFile, "ReverseProxy")

	//convert the string ip into a URL to add headers and cconnecction info
	myURL := parseURL("http://" + serverAddr)
	forwardHeaders(r, myURL)

	//makes the request to the actual server
	response, err := http.DefaultClient.Do(r)
	genericErrHandler(err)

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
	genericErrHandler(err)
	return parsed
}

//Copies headers from the initial request to the forwarded request.
//Also adds the x-forwarded-for header to let the server know the initial request address.
func forwardHeaders(r *http.Request, url *url.URL) {
	r.Host = url.Host
	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.RequestURI = ""

	//gets request's remote address without port number and sets it in forwarding header
	split, _, err := net.SplitHostPort(r.RemoteAddr)
	genericErrHandler(err)
	r.Header.Set("X-Forwarded-For", split)

	//add user and pass for authenticating with the main server
	auth := base64.StdEncoding.EncodeToString([]byte(proxyAuth))
	r.Header.Set("Proxy-Authorization", "Basic "+auth)
}

func genericErrHandler(err error) {
	if err != nil {
		logutil.SendLog(loggerAddr, true, []string{err.Error()}, logFile, "ReverseProxy")
	}
}
