package main

import (
	"architecture/internal/app"
	"architecture/internal/config"
	"log/slog"
	"os"
)

const (
	envlocal = "local"
	envdev   = "dev"
	envprod  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setUpLogger(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env), slog.Any("cfg", cfg), slog.Int("port", cfg.GRPC.Port))

	application := app.New(log, cfg.GRPC.Port)
	application.MustRun()

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
