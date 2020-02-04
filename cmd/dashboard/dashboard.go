package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/mattackard/project-1/pkg/dnsutil"
	"github.com/mattackard/project-1/pkg/logutil"

	procstat "github.com/guillermo/go.procstat"
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
	http.HandleFunc("/stats", serverStats)

	logutil.SendLog(logAddr, false, []string{"Dashboard started at " + noPort}, logFile, "Dashboard")
	http.ListenAndServe(":80", nil)
}

func serverStats(w http.ResponseWriter, r *http.Request) {
	myStats := getStats()
	fmt.Println(myStats)
}

func getStats() *procstat.Stat {
	myStats := procstat.Stat{Pid: os.Getpid()}
	err := myStats.Update()
	if err != nil {
		logutil.SendLog(logAddr, true, []string{err.Error()}, logFile, "Dashboard")
	}
	return &myStats
}
