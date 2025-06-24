package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"pipeline/internal/grpc/handler"
	"pipeline/internal/grpc/registry"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PipelineService struct {
	registryClient registry.RegistryClient
	registryConn   *grpc.ClientConn
	handlerConn    *grpc.ClientConn
}

type HandlerInfo struct {
	Name string
	IP   string
	Port uint32
}

func NewPipelineService(registryAddr string) (*PipelineService, error) {
	// Подключаемся к registry сервису
	registryConn, err := grpc.NewClient(registryAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to registry: %v", err)
	}

	registryClient := registry.NewRegistryClient(registryConn)

	return &PipelineService{
		registryClient: registryClient,
		registryConn:   registryConn,
	}, nil
}

func (s *PipelineService) Close() {
	if s.registryConn != nil {
		s.registryConn.Close()
	}
	if s.handlerConn != nil {
		s.handlerConn.Close()
	}
}

func (s *PipelineService) GetHandlers(ctx context.Context, handlerNames []string) ([]HandlerInfo, error) {
	var handlers []HandlerInfo

	for _, name := range handlerNames {
		req := &registry.GetRequest{
			Name: name,
		}

		resp, err := s.registryClient.Get(ctx, req)
		if err != nil {
			log.Printf("Ошибка получения обработчика %s: %v", name, err)
			continue
		}

		handlers = append(handlers, HandlerInfo{
			Name: name,
			IP:   resp.GetIp(),
			Port: resp.GetPort(),
		})

		log.Printf("Получен обработчик: %s (%s:%d)", name, resp.GetIp(), resp.GetPort())
	}

	return handlers, nil
}

func (s *PipelineService) ProcessMessage(ctx context.Context, message string, handlers []HandlerInfo) (string, error) {
	currentMessage := message

	for i, handlerInfo := range handlers {
		log.Printf("Обработка через %s (%s:%d)", handlerInfo.Name, handlerInfo.IP, handlerInfo.Port)

		handlerAddr := fmt.Sprintf("%s:%d", handlerInfo.IP, handlerInfo.Port)
		handlerConn, err := grpc.NewClient(handlerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return "", fmt.Errorf("failed to connect to handler %s: %v", handlerInfo.Name, err)
		}
		defer handlerConn.Close()

		handlerClient := handler.NewValidateHandlerClient(handlerConn)

		req := &handler.ProcessRequest{
			Message: currentMessage,
		}

		resp, err := handlerClient.Handle(ctx, req)
		if err != nil {
			return "", fmt.Errorf("error processing message in %s: %v", handlerInfo.Name, err)
		}

		currentMessage = resp.GetMessage()
		log.Printf("Результат обработки %s: %s", handlerInfo.Name, currentMessage)

		if i < len(handlers)-1 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	return currentMessage, nil
}
