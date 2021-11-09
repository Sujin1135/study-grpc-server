package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	pb "study-grpc-server/ecommerce"
	"testing"
	"time"
)

func initGRPCServerHTTP2() *grpc.Server {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", port)
	}
	s := grpc.NewServer()
	pb.RegisterProductInfoServer(s, &server{})
	reflection.Register(s)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return s
}

const (
	address = "localhost:50051"
)

func TestServer_AddProduct(t *testing.T) {
	grpcServer := initGRPCServerHTTP2()
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		grpcServer.Stop()
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductInfoClient(conn)

	name := "Samsung S10"
	description := "Samsung Galaxy S10 is the latest smart phone,"
	price := float32(700.0)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	r, err := c.AddProduct(ctx, &pb.Product{Name: name, Description: description, Price: price})
	if err != nil {
		log.Fatalf("Could not add product: %v", err)
	}

	log.Printf("Res %s", r.Value)
	grpcServer.Stop()
}
