package app

import (
	"log/slog"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       uint32
}

func New(log *slog.Logger, regIp string, regPort uint32, port uint32) *App {
	gRPC := grpc.NewServer()
	return &App{log: log, gRPCServer: gRPC, port: port}
}

func (a *App) MustLoad() {
	if err := a.load(); err != nil {
		panic(err.Error())
	}
}

func (a *App) load() error {
	return nil
}
