all:
	go build -ldflags "-w"  -o bin/server server.go
	chmod +x bin/server