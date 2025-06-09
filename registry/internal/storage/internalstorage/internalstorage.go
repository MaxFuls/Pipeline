package internalstorage

import (
	"architecture/internal/domain/models"
	"context"
	"fmt"
)

type Storage struct {
	storage map[string]models.Record
}

func New() *Storage {
	return &Storage{storage: map[string]models.Record{}}
}

func (s *Storage) SaveRecord(ctx context.Context, name string, ip string, port uint32) (err error) {
	_, exist := s.storage[name]

	if exist {
		return fmt.Errorf("record already exists")
	}

	s.storage[name] = models.Record{Name: name, Ip: ip, Port: port}
	return nil
}

func (s *Storage) GetRecord(ctx context.Context, name string) (ip string, port uint32, err error) {
	value, exists := s.storage[name]

	if !exists {
		return "", models.ZeroPort, fmt.Errorf("record does not exists")
	}

	return value.Ip, value.Port, nil
}

func (s *Storage) Contains(ctx context.Context, name string) (contains bool, err error) {
	_, exists := s.storage[name]
	return exists, nil
}
