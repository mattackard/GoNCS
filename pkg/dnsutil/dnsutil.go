package dnsutil

import (
	"fmt"
	"log"
	"net"
)

//Ping send a ping to the DNS so it can record your service and IP
func Ping(address string, serviceName string) {
	conn, err := net.Dial("tcp", address)
	defer conn.Close()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprintf(conn, serviceName)
}

//GetMyIP returns the caller's ip address by sending a blank request to google's DNS server
//and retrieving the local address from the response
func GetMyIP() string {
	conn, err := net.Dial("UDP", "8.8.8.8:80")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	return conn.LocalAddr().String()
}
