package app

import (
	"fmt"
	"log/slog"
	"net"
	"validate_handler/internal/grpc/validate"
	"validate_handler/internal/services"

	"google.golang.org/grpc"
)

type App struct {
	name       string
	ip         string
	port       uint32
	log        *slog.Logger
	gRPCServer *grpc.Server
}

func New(name string, ip string, port uint32, log *slog.Logger, regIp string, regPort uint32) *App {
	gRPC := grpc.NewServer()
	validate.RegisterServer(gRPC, services.New(log))
	return &App{name: name, ip: ip, port: port, log: log, gRPCServer: gRPC}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err.Error())
	}
}

func (a *App) Run() error {
	const op = "App.Run"

	log := a.log.With(slog.String("op", op))

	log.Info("starting " + a.name + " service")

	l, err := net.Listen("tcp", a.ip+fmt.Sprintf(":%d", a.port))

	if err != nil {
		log.Error("Connection failed with error", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	defer l.Close()

	if err := a.gRPCServer.Serve(l); err != nil {
		log.Error("Server error", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) GracefulStop() {
	a.gRPCServer.GracefulStop()
}
