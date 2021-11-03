package main

import (
	"crypto/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	pb "study-grpc-server/ecommerce/ecommerce"
)

const (
	port    = ":50051"
	crtFile = "server.crt"
	keyFile = "server.key"
)

func main() {
	// 1. Generate a certification
	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}

	// 2. Enable the certification
	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}

	// 3. Create a gRPC server instance with the TLS certification
	s := grpc.NewServer(opts...)

	// 4. Register a gRPC server into the created gRPC server
	pb.RegisterProductInfoServer(s, &server{})

	// 5. Create a listener with port(50051)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 6. Bind the gRPC server into the created listener
	if err := s.Serve(lis); err != nil { // 6
		log.Fatalf("failed to serve: %v", err)
	}
}
