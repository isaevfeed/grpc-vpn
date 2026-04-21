package main

import (
	"context"
	"grpc-vpn/internal/transport/grpc"
	"grpc-vpn/internal/transport/socks5"
	"log"
	"log/slog"
	"net"
	"os"
	"time"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	grpcClient, err := grpc.NewClient("localhost:50051")
	if err != nil {
		log.Println(err)
	}

	srv := socks5.New(":1080", slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	srv.Start(func(conn net.Conn, target string) {
		err = grpcClient.Handle(ctx, conn, target)
		if err != nil {
			log.Println(err)
		}
	})
}
