package main

import (
	"context"
	"log/slog"
	"os"
	"strconv"
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
	log.Info("starting pipeline service", slog.String("env", cfg.Env), slog.String("registry addres", cfg.Registry.Ip+":"+strconv.Itoa(int(cfg.GRPC.Port))))

	ctx := context.Background()

	port := fmt.Sprint(cfg.Registry.Port)
	addr := cfg.Registry.Ip + ":" + port

	conn, _ := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	defer conn.Close()

	client := registry.NewRegistryClient(conn)

	client.Register(ctx, &registry.RegistryRequest{Name: "validation", Ip: "192.168.1.0", Port: 5001})

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
