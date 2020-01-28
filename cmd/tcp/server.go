package main

import (
	"log"
	"net"
	"os"
)

var port = os.Getenv("TCPPORT")

func main() {
	l, err := net.Listen("tcp", port)
	log.Printf("tcp server is running on port %s\n", port)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	//wait for a connection
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(c net.Conn) {
			c.Write([]byte("HTTP/1.1 200 OK"))
			c.Close()
		}(conn)
	}
}
