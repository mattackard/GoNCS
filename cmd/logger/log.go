package logger

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

//SendLog send the log message over tcp and throws an error if the log message is an error
func SendLog(address string, isErr bool, data []string) {
	conn, err := net.Dial("tcp", address+logPort)
	defer conn.Close()
	if err != nil {
		log.Fatalln(err)
	}
	for _, v := range data {
		conn.Write([]byte(v))
	}
	if isErr {
		log.Fatalln(data)
	}
}

//LogServerRequest creates a summary of the http connection information and send it to the connected logger
func LogServerRequest(w http.ResponseWriter, r *http.Request, loggerAddr string) {
	method := r.Method
	url := r.URL
	httpVer := r.Proto
	host := r.Host
	closeConn := r.Close
	address := r.RemoteAddr
	reqData := fmt.Sprint(method, url, httpVer, host, closeConn, address)
	SendLog(loggerAddr, false, []string{reqData})
}
