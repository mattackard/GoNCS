//Package perfutil provides functions and structs for storing and transporting performace data
package perfutil

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/segmentio/stats/procstats"
)

//ContainerStats holds the array of service structs containing service runtime stats
type ContainerStats struct {
	Containers []Service `json:"containers"`
}

//Service contains the runtime stats for a process
type Service struct {
	ServiceName  string        `json:"serviceName"`
	CPUShare     int64         `json:"cpuShare"`
	CPUUserTime  time.Duration `json:"cpuUserTime"`
	CPUSysTime   time.Duration `json:"cpuSysTime"`
	AvailableMem uint64        `json:"availableMem"`
	MemUsage     uint64        `json:"memUsage"`
	OpenFiles    uint64        `json:"openFiles"`
	ThreadCount  uint64        `json:"threadCount"`
	CPUDelay     time.Duration `json:"cpuDelay"`
	BlockIODelay time.Duration `json:"blockIODelay"`
}

//GetServerStats returns a Serivce struct containing the stats of the container calling the function
func GetServerStats() Service {
	myStats, err := procstats.CollectProcInfo(os.Getpid())
	myDelay, err := procstats.CollectDelayInfo(os.Getegid())
	if err != nil {
		log.Fatalln(err)
	}
	dashboardStats := Service{
		"Dashboard",
		myStats.CPU.Shares,
		myStats.CPU.User,
		myStats.CPU.Sys,
		myStats.Memory.Available,
		myStats.Memory.Size,
		myStats.Files.Open,
		myStats.Threads.Num,
		myDelay.CPUDelay,
		myDelay.BlockIODelay,
	}
	return dashboardStats
}

//RequestStatsHTTP makes an HTTP request to requestAddr/getStats and puts the response into a service struct
func RequestStatsHTTP(requestAddr string) Service {
	resp, err := http.Get(requestAddr + "/getStats")
	buffer := make([]byte, 1024)
	var stats Service
	_, err = resp.Body.Read(buffer)
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(buffer, &stats)
	return stats
}
