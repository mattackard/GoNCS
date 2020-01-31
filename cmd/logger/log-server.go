package main

import (
	"log"
	"os"

	"github.com/mattackard/project-1/pkg/logger"
	//"github.com/mattackard/project-1/pkg/logger"
)

var logPort = os.Getenv("LOGPORT")
var logAddress = os.Getenv("LOGGERNAME")
var fullAddress = logAddress + ":" + logPort

func main() {
	logFile, err := os.Create("./logs/logs.txt")
	if err != nil {
		log.Fatalln(err)
	}
	logger.CreateLogServerAndListen(fullAddress, logFile)
}
