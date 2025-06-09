package main

import (
	"log/slog"
	"os"
	"pipeline/internal/config"
	"strconv"
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
