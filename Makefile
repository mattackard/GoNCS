include .env

proxy: 
	go run ./cmd/reverseproxy/proxy.go

tcpserver:
	go run ./cmd/tcp/server.go