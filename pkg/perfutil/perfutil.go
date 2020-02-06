//Package perfutil provides functions and structs for storing and transporting performace data
package perfutil

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/mattackard/project-1/pkg/dnsutil"

	"github.com/segmentio/stats/procstats"
)

//ContainerStats holds the array of service structs containing service runtime stats
type ContainerStats struct {
	Containers []Service `json:"containers"`
}

//Service contains the runtime stats for a process
type Service struct {
	ServiceName  string        `json:"serviceName"`
	IP           string        `json:"ip"`
	CPUShare     int64         `json:"cpuShare"`
	CPUUserTime  time.Duration `json:"cpuUserTime"`
	CPUSysTime   time.Duration `json:"cpuSysTime"`
	AvailableMem uint64        `json:"availableMem"`
	MemUsage     uint64        `json:"memUsage"`
	OpenFiles    uint64        `json:"openFiles"`
	ThreadCount  uint64        `json:"threadCount"`
}

//GetServerStats returns a Serivce struct containing the stats of the container calling the function
func GetServerStats() Service {
	myStats, err := procstats.CollectProcInfo(os.Getpid())
	if err != nil {
		log.Fatalln(err)
	}
	containerStats := Service{
		"Dashboard",
		dnsutil.TrimPort(dnsutil.GetMyIP()),
		myStats.CPU.Shares,
		myStats.CPU.User,
		myStats.CPU.Sys,
		myStats.Memory.Available,
		myStats.Memory.Size,
		myStats.Files.Open,
		myStats.Threads.Num,
	}
	return containerStats
}

//RequestStatsHTTP makes an HTTP request to requestAddr/getStats and puts the response into a service struct
func RequestStatsHTTP(requestAddr string) Service {
	resp, err := http.Get("http://" + requestAddr + "/getStats")
	if err != nil {
		log.Fatalln(err)
	}
	var stats Service

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(body, &stats)
	if err != nil {
		log.Fatalln(err)
	}
	return stats
}

//RequestStatsTCP makes a TCP call, asks for container stats, and records the response into a Service struct
func RequestStatsTCP(requestAddr string) Service {
	var conn net.Conn
	var err error
	var stats Service
	for {
		conn, err = net.Dial("tcp", requestAddr)
		if err == nil {
			break
		}
	}
	defer conn.Close()
	conn.Write([]byte("containerStats"))

	buffer := make([]byte, 1024)
	conn.Read(buffer)

	//trim any extra nil bytes
	buffer = bytes.Trim(buffer, "\x00")

	err = json.Unmarshal(buffer, &stats)
	if err != nil {
		log.Fatalln(err)
	}
	return stats
}

//SendStatsHTTP is a generic http handler for getting the server's stats and sending them in the response
func SendStatsHTTP(w http.ResponseWriter, r *http.Request) {
	stats := GetServerStats()
	response, err := json.Marshal(stats)
	if err != nil {
		log.Fatalln(err)
	}
	w.Write(response)
}

//SendStatsTCP is a generic tcp handler for getting the server's stats and sending them back through the connection
func SendStatsTCP(conn net.Conn) {
	stats := GetServerStats()
	response, err := json.Marshal(stats)
	if err != nil {
		log.Fatalln(err)
	}
	conn.Write(response)
	conn.Close()
}
