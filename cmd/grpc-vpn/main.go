package main

import (
	"grpc-vpn/internal/app"
	"log/slog"
	"os"
)

func main() {
	//TODO сконфигурировать логгер
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app.New(log).MustRun()
}
