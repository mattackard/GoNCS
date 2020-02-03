package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/mattackard/project-1/pkg/dnsutil"
	"github.com/mattackard/project-1/pkg/logutil"
)

//DNS holds the name and ip address of each service that connects to it
var DNS = make(map[string]string)
var dnsPort = os.Getenv("DNSPORT")

var logPort = os.Getenv("LOGPORT")
var logName = os.Getenv("LOGGERNAME")
var loggerAddr = logName + ":" + logPort

func main() {

	logFile := logutil.OpenLogFile("./logs/")
	defer logFile.Close()

	l, err := net.Listen("tcp", ":"+dnsPort)
	if err != nil {
		logutil.SendLog(loggerAddr, true, []string{err.Error()}, logFile, "DNS")
	}
	defer l.Close()

	//send messages to log file to record startup
	dnsIP := dnsutil.GetMyIP()
	logutil.SendLog(loggerAddr, false, []string{"DNS Server started at " + dnsIP.String()}, logFile, "DNS")

	//wait for a connection
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		buffer := make([]byte, 1024)
		conn.Read(buffer)

		//read the service name sent and assign it using it's IP in the dns's map
		bufferText := string(bytes.Trim(buffer, "\x00"))
		fmt.Println(bufferText, conn.RemoteAddr().String())
		DNS[bufferText] = conn.RemoteAddr().String()
		logutil.WriteToLog(logFile, "DNS", []string{bufferText + " started at " + DNS[bufferText]})
		go func(c net.Conn) {
			c.Write(buffer)
			c.Close()
		}(conn)
	}
}
