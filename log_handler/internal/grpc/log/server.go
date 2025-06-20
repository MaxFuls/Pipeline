package log

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type ServerAPI struct {
	UnimplementedValidateHandlerServer
	validatation ValidateHandler
}

type ValidateHandler interface {
	Handle(ctx context.Context, message string) (string, error)
}

func RegisterServer(grpc *grpc.Server, handler ValidateHandler) {
	RegisterValidateHandlerServer(grpc, &ServerAPI{validatation: handler})
}

func (s *ServerAPI) Handle(ctx context.Context, req *ProcessRequest) (*ProcessResponse, error) {
	message := req.GetMessage()
	handled, err := s.validatation.Handle(ctx, message)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ProcessResponse{Message: handled}, nil
}
