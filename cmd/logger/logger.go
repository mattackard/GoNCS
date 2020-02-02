package main

import (
	"os"

	"github.com/mattackard/project-1/pkg/dnsutil"
	"github.com/mattackard/project-1/pkg/logutil"
)

var logPort = os.Getenv("LOGPORT")
var logAddress = os.Getenv("LOGGERNAME")
var dnsPort = os.Getenv("DNSPORT")
var dnsName = os.Getenv("DNSNAME")

func main() {
	logFile := logutil.OpenLogFile("./logs")
	defer logFile.Close()

	//send messages to log file and dns to record startup
	loggerIP := dnsutil.Ping(dnsName+":"+dnsPort, "logger")
	logutil.WriteToLog(logFile, []string{"Reverse Proxy Server started at " + loggerIP})

	logutil.CreateLogServerAndListen(logAddress, logPort, logFile)
}
