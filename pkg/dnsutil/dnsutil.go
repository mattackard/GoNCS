package dnsutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
)

//DNS holds the addresses of all currently running services as an array of [name, address]
type DNS struct {
	Services map[string]string `json:"services"`
}

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
	fmt.Fprintf(conn, "recordAddress="+serviceName)
	ip := conn.LocalAddr()
	return ip
}

//GetServiceAddresses return an array of string
func GetServiceAddresses(dnsAddr string) DNS {
	var conn net.Conn
	var err error

	//wait for a successful connection
	for {
		conn, err = net.Dial("tcp", dnsAddr)
		if err == nil {
			break
		}
	}
	defer conn.Close()

	//make a request for all DNS addresses and place response into a buffer
	conn.Write([]byte("getAllAddresses"))
	buffer := make([]byte, 1024)
	conn.Read(buffer)

	//trim any extra nil bytes
	buffer = bytes.Trim(buffer, "\x00")

	//unmarshal the json back into a DNS struct and return
	var response DNS
	err = json.Unmarshal(buffer, &response)
	if err != nil {
		log.Fatal(err)
	}
	return response
}

//GetMyIP returns the caller's ip address by sending a blank request to google's DNS server
//and retrieving the local address from the response
func GetMyIP() net.Addr {
	var conn net.Conn
	var err error
	for {
		conn, err = net.Dial("udp", "8.8.8.8:80")
		if err == nil {
			break
		}
	}
	defer conn.Close()
	ip := conn.LocalAddr()
	return ip
}

//GetServiceIP queries the dns server for a currently running service and returns the IP address
func GetServiceIP(dnsAddr string, serviceName string) string {
	var conn net.Conn
	var err error
	for {
		conn, err = net.Dial("tcp", dnsAddr)
		if err == nil {
			break
		}
	}
	defer conn.Close()
	fmt.Fprint(conn, "getAddress="+serviceName)

	buffer := make([]byte, 1024)
	conn.Read(buffer)

	//trim te null characters from the buffer and convert to string
	bufferText := string(bytes.Trim(buffer, "\x00"))

	//remove the "servicename=" and leave just the IP
	bufferText = strings.Split(bufferText, "=")[1]
	return bufferText
}

//TrimPort converts an IP address into a string containing the IP without the port attached
func TrimPort(ip net.Addr) string {
	stringIP := ip.String()
	noPort := strings.Split(stringIP, ":")[0]
	return noPort
}
