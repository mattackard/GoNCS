package main

import (
	"net/http"
	"os"

	"github.com/mattackard/project-1/pkg/dnsutil"
	"github.com/mattackard/project-1/pkg/logutil"
)

var dnsAddr string
var logAddr = "logger:6060"
var logFile *os.File

func main() {
	logFile = logutil.OpenLogFile("/logs")

	myIP := dnsutil.Ping("dns:6060", "dashboard")
	noPort := dnsutil.TrimPort(myIP)

	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/dashboard/", http.StripPrefix("/dashboard/", fs))

	logutil.SendLog(logAddr, false, []string{"Dashboard started at " + noPort}, logFile, "Dashboard")
	http.ListenAndServe(":80", nil)
}
