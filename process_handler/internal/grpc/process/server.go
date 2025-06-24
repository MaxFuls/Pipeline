package process

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type ServerAPI struct {
	UnimplementedValidateHandlerServer
	processor ProcessHandler
}

type ProcessHandler interface {
	Handle(ctx context.Context, message string) (string, error)
}

func RegisterServer(grpc *grpc.Server, handler ProcessHandler) {
	RegisterValidateHandlerServer(grpc, &ServerAPI{processor: handler})
}

func (s *ServerAPI) Handle(ctx context.Context, req *ProcessRequest) (*ProcessResponse, error) {
	message := req.GetMessage()
	handled, err := s.processor.Handle(ctx, message)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ProcessResponse{Message: handled}, nil
}
