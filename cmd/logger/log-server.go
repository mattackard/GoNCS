package main

import (
	"os"

	"github.com/mattackard/project-1/pkg/logger"
)

var logPort = os.Getenv("LOGPORT")
var logAddress = os.Getenv("LOGGERNAME")

func main() {
	logger.CreateLogServerAndListen(logAddress, logPort)
}
