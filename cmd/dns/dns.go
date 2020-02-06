package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"

	"github.com/mattackard/project-1/pkg/dnsutil"
	"github.com/mattackard/project-1/pkg/logutil"
	"github.com/mattackard/project-1/pkg/perfutil"
)

//DNS holds the name and ip address of each service that connects to it
var DNS = dnsutil.DNS{}

var universalPort = "6060"

var logName = os.Getenv("LOGGERNAME")
var loggerAddr = logName + ":" + universalPort

func main() {

	logFile := logutil.OpenLogFile("./logs/")
	defer logFile.Close()

	//initialize services map and add dns to it
	DNS.Services = make(map[string]string)
	dnsIP := dnsutil.GetMyIP()
	DNS.Services["DNS"] = dnsutil.TrimPort(dnsIP) + ":" + universalPort

	//create tcp server
	l, err := net.Listen("tcp", ":"+universalPort)
	if err != nil {
		logutil.SendLog(loggerAddr, true, []string{err.Error()}, logFile, "DNS")
	}
	defer l.Close()

	//send messages to log file to record startup
	logutil.SendLog(loggerAddr, false, []string{"DNS Server started at " + DNS.Services["DNS"]}, logFile, "DNS")

	//wait for a connection
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		//set up a buffer to read and write to the tcp connection
		buffer := make([]byte, 1024)
		conn.Read(buffer)

		//trim the nil bytes from buffer and split to locate subcommands
		bufferText := string(bytes.Trim(buffer, "\x00"))
		bufferSlice := strings.Split(bufferText, "=")

		//check subcommands sent through tcp
		if bufferSlice[0] == "recordAddress" {

			//read the service name and assign it in the DNS map to the remote address
			DNS.Services[bufferSlice[1]] = dnsutil.TrimPort(conn.RemoteAddr()) + ":" + universalPort

			//record the connection registration to the logger and responsd with service name
			logutil.SendLog(loggerAddr, false, []string{bufferSlice[1] + " started at " + DNS.Services[bufferSlice[1]]}, logFile, "DNS")
			conn.Write([]byte(bufferSlice[1]))

		} else if bufferSlice[0] == "getAddress" {
			//send the service name and ip:port in the response
			logutil.SendLog(loggerAddr, false, []string{dnsutil.TrimPort(conn.RemoteAddr()) + " requested the address for " + bufferSlice[1]}, logFile, "DNS")
			conn.Write([]byte(bufferSlice[1] + "=" + DNS.Services[bufferSlice[1]]))

		} else if bufferSlice[0] == "getAllAddresses" {
			//send the whole DNS map
			logutil.WriteToLog(logFile, "DNS", []string{dnsutil.TrimPort(conn.RemoteAddr()) + " requested all DNS addresses"})

			//encode the DNS struct into a json object
			bytes, err := json.Marshal(DNS)
			if err != nil {
				logutil.SendLog(loggerAddr, true, []string{err.Error()}, logFile, "DNS")
			}
			conn.Write(bytes)

		} else if bufferSlice[0] == "containerStats" {
			//send the dns process stats
			perfutil.SendStatsTCP(conn)

		} else {
			//don't allow data without a subcommand
			conn.Write([]byte("400 Bad Request"))
		}
	}
}
