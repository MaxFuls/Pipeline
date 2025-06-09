package registry

import (
	"architecture/internal/domain/models"
	"context"
	"log/slog"
)

type Registry struct {
	log      *slog.Logger
	saver    RecordsSaver
	provider RecordsProvider
}

func (reg *Registry) Register(ctx context.Context, name string, ip string, port uint32) (string, error) {
	contains, err := reg.provider.Contains(ctx, name)

	if err != nil {
		return "", err
	}

	if contains {
		return "record already exists", nil
	}

	err = reg.saver.SaveRecord(ctx, name, ip, port)

	if err != nil {
		return "", err
	}
	return "record was succesfuly added", nil
}

func (reg *Registry) Get(ctx context.Context, name string) (string, string, uint32, error) {

	ip, port, err := reg.provider.GetRecord(ctx, name)

	if err != nil {
		return "", "", models.ZeroPort, err
	}

	return "record is found", ip, port, nil
}

type RecordsSaver interface {
	SaveRecord(ctx context.Context, name string, ip string, port uint32) (err error)
}

type RecordsProvider interface {
	GetRecord(ctx context.Context, name string) (ip string, port uint32, err error)
	Contains(ctx context.Context, name string) (contains bool, err error)
}

func New(log *slog.Logger, saver RecordsSaver, provider RecordsProvider) *Registry {
	return &Registry{
		log:      log,
		saver:    saver,
		provider: provider,
	}
}
