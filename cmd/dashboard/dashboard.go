package main

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/mattackard/project-1/pkg/dnsutil"
	"github.com/mattackard/project-1/pkg/logutil"

	"github.com/segmentio/stats/procstats"
)

//ContainerStats holds the array of service structs containing service runtime stats
type ContainerStats struct {
	Containers []Service `json:"containers"`
}

//Service contains the runtime stats for a process
type Service struct {
	ServiceName string
	CPUShare    string
	CPUUserTime string
	CPUSysTime  string
}

var dnsAddr string
var logAddr = "logger:6060"
var logFile *os.File

func main() {
	logFile = logutil.OpenLogFile("/logs")

	myIP := dnsutil.Ping("dns:6060", "dashboard")
	noPort := dnsutil.TrimPort(myIP)

	//file server for html file
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/dashboard/", http.StripPrefix("/dashboard/", fs))

	//endpoints for javascript requests
	http.HandleFunc("/stats", serverStats)
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

func serverStats(w http.ResponseWriter, r *http.Request) {
	w = setHeaders(w)
	myStats, err := procstats.CollectProcInfo(os.Getpid())
	if err != nil {
		logutil.SendLog(logAddr, true, []string{err.Error()}, logFile, "Dashboard")
	}
	fmt.Fprint(w, myStats)
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
