package logger

import (
	"fmt"
	"log"
	"net"
	"os"
)

var logPort = os.Getenv("LOGPORT")
var logAddress = os.Getenv("DOCKERLOGNAME")

func main() {
	CreateLogServerAndListen(logAddress, logPort)
}

//CreateLogServerAndListen runs a tcp server at address:port
func CreateLogServerAndListen(address string, port string) {
	l, err := net.Listen("tcp", address+": "+port)
	log.Printf("Logger is listening on port %s\n", port)
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
		buffer := make([]byte, 1024)
		conn.Read(buffer)
		fmt.Println(string(buffer))
		go func(c net.Conn) {
			c.Write(buffer)
			c.Close()
		}(conn)
	}
}
