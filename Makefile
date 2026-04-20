run-vpn:
	go run cmd/grpc-vpn/main.go

build-vpn:
	go build -o bin/grpc-vpn cmd/grpc-vpn/main.go
