package main

import (
	"grpc-vpn/internal/app"
	"grpc-vpn/internal/transport/grpc"
	"log/slog"
	"os"
)

func main() {
	//TODO сконфигурировать логгер
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	grpcSrv := grpc.New()

	app.New(log, grpcSrv).MustRun()
}
