package main

import (
	"context"
	"fmt"
	"pipeline/internal/grpc/registry"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, _ := grpc.NewClient("localhost:5000", grpc.WithTransportCredentials(insecure.NewCredentials()))

	defer conn.Close()

	client := registry.NewRegistryClient(conn)

	ctx := context.Background()

	response, _ := client.Register(ctx, &registry.RegistryRequest{Name: "test", Ip: "124.111.1.1", Port: 123})

	fmt.Print(response)
}
