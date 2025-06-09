package registry

import (
	"context"
	"net"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Registry interface {
	Register(ctx context.Context, name string, ip string, port uint32) (string, error)
	Get(ctx context.Context, name string) (string, string, uint32, error)
}

type serverAPI struct {
	UnimplementedRegistryServer
	registry Registry
}

func RegisterServer(gRPC *grpc.Server, registry Registry) {
	RegisterRegistryServer(gRPC, &serverAPI{registry: registry})
}

func (s *serverAPI) Register(ctx context.Context, req *RegistryRequest) (*RegistryResponse, error) {
	name := req.GetName()
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name of service is required")
	}

	ip := req.GetIp()

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil || parsedIP.To4() == nil {
		return nil, status.Error(codes.InvalidArgument, "ip is not a valid ipv4 address")
	}

	port := req.GetPort()

	if port <= 0 || port > 65535 {
		return nil, status.Error(codes.InvalidArgument, "port must be between 1 and 65535")
	}

	msg, err := s.registry.Register(ctx, name, ip, port)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &RegistryResponse{Message: msg}, nil
}

func (s *serverAPI) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	name := req.GetName()
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	msg, ip, port, err := s.registry.Get(ctx, name)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &GetResponse{Message: msg, Ip: ip, Port: port}, nil
}
