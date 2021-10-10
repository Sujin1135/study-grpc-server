package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	pb "study-grpc-server/ecommerce/github.com/Sujin1135/study-grpc-server/blob/master/ecommerce"
	pb2 "study-grpc-server/order/order"
)

const (
	port = ":50051"
)

func main() {
	initSampleData()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProductInfoServer(s, &server{})
	pb2.RegisterOrderManagementServer(s, &orderServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
