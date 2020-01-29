include .env

proxy: 
	go run ./cmd/reverseproxy/proxy.go

tcpserver:
	go run ./cmd/tcp/server.go

all:
	docker build -t rproxy ./cmd/reverseproxy/
	docker build -t noteserver $GOPATH/src/github.com/mattackard/project-0/cmd/GoNotesd/
	docker-compose up