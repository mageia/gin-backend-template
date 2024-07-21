package service

import (
	"api-server/rpc"
	"context"
)

type HelloWorld struct{}

func (h *HelloWorld) SayHello(ctx context.Context, req *rpc.HelloRequest) (*rpc.HelloResponse, error) {
	return &rpc.HelloResponse{Message: "Hello " + req.Name}, nil
}
