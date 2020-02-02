package dnsutil

import (
	"fmt"
	"log"
	"net"
)

//Ping send a ping to the DNS so it can record your service and IP
func Ping(address string, serviceName string) net.Addr {
	var conn net.Conn
	var err error
	for {
		conn, err = net.Dial("tcp", address)
		if err == nil {
			break
		}
	}
	defer conn.Close()
	fmt.Fprintf(conn, serviceName)
	ip := conn.LocalAddr()
	return ip
}

//GetMyIP returns the caller's ip address by sending a blank request to google's DNS server
//and retrieving the local address from the response
func GetMyIP() net.Addr {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	ip := conn.LocalAddr()
	return ip
}
