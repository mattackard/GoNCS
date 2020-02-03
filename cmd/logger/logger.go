package main

import (
	"os"

	"github.com/mattackard/project-1/pkg/dnsutil"
	"github.com/mattackard/project-1/pkg/logutil"
)

var dnsPort = os.Getenv("DNSPORT")
var dnsName = os.Getenv("DNSNAME")
var dnsAddr string

var logPort = os.Getenv("LOGPORT")
var logName = os.Getenv("LOGGERNAME")

func init() {
	//set the default dns address if none are specified in the environment
	if dnsName == "" || dnsPort == "" {
		dnsAddr = "dns:6060"
	} else {
		dnsAddr = dnsName + ":" + dnsPort
	}

	//same as above but for logger address
	if logName == "" || logPort == "" {
		logName = "logger"
		logPort = "6060"
	}
}

func main() {
	logFile := logutil.OpenLogFile("./logs")
	defer logFile.Close()

	//send messages to log file and dns to record startup
	dnsutil.Ping(dnsAddr, logName)

	logutil.CreateLogServerAndListen(logName, logPort, logFile)
}
