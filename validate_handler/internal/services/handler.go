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
	if len(message)%2 == 0 {
		message = message + " - accepted"
	} else {
		message = message + " - rejected"
	}
	return message, nil
}
