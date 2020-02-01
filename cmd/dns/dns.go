package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/mattackard/project-1/pkg/logutil"
)

//DNS holds the name and ip address of each service that connects to it
var DNS = make(map[string]string)
var dnsPort = os.Getenv("DNSPORT")

var logPort = os.Getenv("LOGPORT")
var logName = os.Getenv("LOGGERNAME")
var loggerAddr = logName + ":" + logPort
var logFile *os.File

func init() {
	logFile = logutil.OpenLogFile("./logs/")
}

func main() {
	l, err := net.Listen("tcp", dnsPort)
	log.Printf("DNS is listening at %s\n", dnsPort)
	if err != nil {
		logutil.SendLog(loggerAddr, true, []string{err.Error()}, logFile, "DNS")
	}
	defer l.Close()

	//wait for a connection
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		buffer := make([]byte, 1024)
		conn.Read(buffer)
		bufferText := string(buffer) + "\n"
		serviceInfo := strings.Split(bufferText, ":")
		fmt.Print(serviceInfo)
		DNS[serviceInfo[0]] = serviceInfo[1]
		go func(c net.Conn) {
			c.Write(buffer)
			c.Close()
		}(conn)
	}
}
