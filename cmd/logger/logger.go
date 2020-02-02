package main

import (
	"os"

	"github.com/mattackard/project-1/pkg/logutil"
)

var logPort = os.Getenv("LOGPORT")
var logAddress = os.Getenv("LOGGERNAME")
var fullAddress = logAddress + ":" + logPort

func main() {
	logFile := logutil.OpenLogFile("./logs")
	defer logFile.Close()
	logutil.CreateLogServerAndListen(fullAddress, logFile)
}
