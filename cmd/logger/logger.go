package main

import (
	"os"

	"github.com/mattackard/project-1/pkg/logutil"
)

var logPort = os.Getenv("LOGPORT")
var logAddress = os.Getenv("LOGGERNAME")

func main() {
	logFile := logutil.OpenLogFile("./logs")
	defer logFile.Close()
	logutil.CreateLogServerAndListen(logAddress, logPort, logFile)
}
