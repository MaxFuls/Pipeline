package app

import (
	"architecture/internal/grpc/registry"
	"architecture/internal/storage/internalstorage"
	"fmt"
	"log/slog"
	"net"

	reg "architecture/internal/services/registry"

	"google.golang.org/grpc"
)

type App struct {
	log         *slog.Logger
	gRPCSerever *grpc.Server
	port        int
}

func New(log *slog.Logger, port int) *App {
	gRPCServer := grpc.NewServer()
	storage := internalstorage.New()
	registry.RegisterServer(gRPCServer, reg.New(log, storage, storage))
	return &App{log: log, port: port, gRPCSerever: gRPCServer}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "app.Run"

	log := a.log.With(slog.String("op", op))

	log.Info("starting application")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))

	if err != nil {
		log.Error("Connection failed with error", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	defer l.Close()

	if err := a.gRPCSerever.Serve(l); err != nil {
		log.Error("Server error", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
