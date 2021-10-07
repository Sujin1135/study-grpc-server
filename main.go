package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	pb "study-grpc-server/ecommerce/github.com/Sujin1135/study-grpc-server/blob/master/ecommerce"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProductInfoServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
