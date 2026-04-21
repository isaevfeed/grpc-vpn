package app

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type (
	App struct {
		log *slog.Logger

		grpcSrv *grpc.Server
	}
)

func New(log *slog.Logger, grpcSrv *grpc.Server) *App {
	return &App{
		log:     log,
		grpcSrv: grpcSrv,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		a.log.Error("failed to run", err)
		panic(err.Error())
	}
}

func (a *App) Run() error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	a.log.Info("grpc server started on 50051")

	err = a.grpcSrv.Serve(lis)
	if err != nil {
		return fmt.Errorf("grpcSrv.Serve: %w", err)
	}

	return nil
}
