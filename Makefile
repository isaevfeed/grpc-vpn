run-vpn-server:
	go run cmd/grpc-vpn/main.go

build-vpn-server:
	go build -o bin/grpc-vpn cmd/grpc-vpn/main.go

run-vpn-client:
	go run cmd/grpc-vpn-client/main.go

build-vpn-client:
	go build -o bin/grpc-vpn cmd/grpc-vpn-client/main.go
