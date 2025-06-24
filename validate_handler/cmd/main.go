package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"validate_handler/internal/app"
	"validate_handler/internal/config"
	"validate_handler/internal/grpc/registry"

	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	envlocal = "local"
	envdev   = "dev"
	envprod  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setUpLogger(cfg.Env)
	log.Info("starting "+cfg.Name+" service", slog.String("env", cfg.Env), slog.String("registry addres ", cfg.Registry.Ip+":"+strconv.Itoa(int(cfg.GRPC.Port))))

	app := app.New(cfg.Name, cfg.GRPC.Ip, cfg.GRPC.Port, log, cfg.Registry.Ip, cfg.Registry.Port)

	log.Info("trying to run service")
	go app.MustRun()
	log.Info("application is running")

	log.Info("getting context")
	ctx := context.Background()

	port := fmt.Sprint(cfg.Registry.Port)
	addr := cfg.Registry.Ip + ":" + port

	log.Info("try to connect to registry service")
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Error("connection failed", slog.String("err", err.Error()))
		return
	}

	log.Info("connection is established")
	defer conn.Close()

	client := registry.NewRegistryClient(conn)

	resp, err := client.Register(ctx, &registry.RegistryRequest{Name: cfg.Name, Ip: cfg.GRPC.Ip, Port: cfg.GRPC.Port})

	if err != nil {
		log.Error("registry failed", slog.String("error", err.Error()))
	} else {
		log.Info("service was succesfully registrated with response", slog.String("response", resp.GetMessage()))
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Info("Recieved SIGINT or SIGTERM signal, stopping")
	app.GracefulStop()

}

func setUpLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envlocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envdev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envprod:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
