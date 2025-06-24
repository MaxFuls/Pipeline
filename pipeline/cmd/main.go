package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"pipeline/api/rest"
	"pipeline/internal/app"
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
	log.Info("starting pipeline service", slog.String("env", cfg.Env), slog.String("registry addres", cfg.Registry.Ip+":"+strconv.Itoa(int(cfg.Registry.Port))))

	registryAddr := fmt.Sprintf("%s:%d", cfg.Registry.Ip, cfg.Registry.Port)
	pipelineService, err := app.NewPipelineService(registryAddr)
	if err != nil {
		log.Error("failed to create pipeline service", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pipelineService.Close()

	httpApi := rest.New(pipelineService)

	http.HandleFunc("/", httpApi.HandleHome)
	http.HandleFunc("/submit", httpApi.HandleSubmit)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := cfg.GRPC.Port
	addr := fmt.Sprintf(":%d", port)

	log.Info("starting HTTP server", slog.String("address", addr))
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Error("failed to start HTTP server", slog.String("error", err.Error()))
		os.Exit(1)
	}
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
