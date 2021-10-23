package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
)

func orderUnaryInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("===== [Server Interceptor] in previous process", info.FullMethod)
	log.Printf("In previous process")

	m, err := handler(ctx, req)

	log.Printf("Post Proc Message : %s", m)

	return m, err
}
