include .env

proxy: 
	go run ./cmd/reverseproxy/proxy.go

tcpserver:
	go run ./cmd/tcp/server.go

all:
	go build ../project-0/cmd/GoNotesClient/client.go 
	./client &
	docker-compose up
	rm client