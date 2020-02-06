package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/mattackard/project-1/pkg/dnsutil"
	"github.com/mattackard/project-1/pkg/logutil"
	"github.com/mattackard/project-1/pkg/perfutil"
)

var dnsAddr = "dns:6060"
var logAddr = "logger:6060"
var logFile *os.File

func main() {
	logFile = logutil.OpenLogFile("/logs")
	logAddr = dnsutil.GetServiceIP(dnsAddr, "logger")

	myIP := dnsutil.Ping(dnsAddr, "dashboard")
	noPort := dnsutil.TrimPort(myIP)

	//file server for html file
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/dashboard/", http.StripPrefix("/dashboard/", fs))

	//endpoints for javascript requests
	http.HandleFunc("/stats", getAllStats)
	http.HandleFunc("/getLogs", getMasterLog)

	logutil.SendLog(logAddr, false, []string{"Dashboard started at " + noPort}, logFile, "Dashboard")
	http.ListenAndServe(":80", nil)
}

func setHeaders(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	return w
}

// //gets all stats from all online services and send the data to the dashboard client
func getAllStats(w http.ResponseWriter, r *http.Request) {

	//create container instance for stats
	var containerStats perfutil.ContainerStats

	//get all currently running container IPs from dns
	addresses := dnsutil.GetServiceAddresses(dnsAddr)
	fmt.Println(addresses)

	//for each address, send a request for stats and append to containerStats
	for _, v := range addresses.Services {
		fmt.Println(v)
		thisService := perfutil.RequestStatsHTTP(v)
		containerStats.Containers = append(containerStats.Containers, thisService)
	}

	//append the local container stats for the dashboard last
	containerStats.Containers = append(containerStats.Containers, perfutil.GetServerStats())

	//marshal containerStats into a byte stream to be recieved by the client as json
	bytes, err := json.Marshal(containerStats)
	if err != nil {
		logutil.SendLog(logAddr, true, []string{err.Error()}, logFile, "Dashboard")
	}

	w = setHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func getMasterLog(w http.ResponseWriter, r *http.Request) {
	var conn net.Conn
	var err error
	for {
		conn, err = net.Dial("tcp", logAddr)
		if err == nil {
			break
		}
	}
	defer conn.Close()

	//request the master log file's contents
	conn.Write([]byte("sendLog"))

	//create a buffer and read the response from it
	buffer := make([]byte, 16384)
	var bufferText string
	conn.Read(buffer)

	//trim any extra nil bytes
	bufferText = string(bytes.Trim(buffer, "\x00"))

	w = setHeaders(w)
	fmt.Fprint(w, bufferText)
}
