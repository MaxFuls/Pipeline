package services

import (
	"context"
	"log/slog"
)

type ValidateHandler struct {
	log *slog.Logger
}

func New(log *slog.Logger) *ValidateHandler {
	return &ValidateHandler{log: log}
}

func (h *ValidateHandler) Handle(ctx context.Context, message string) (string, error) {
	return "", nil
}
