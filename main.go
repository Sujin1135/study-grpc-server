package study_grpc_server

import (
	"context"
	"google.golang.org/grpc"
	pb "github.com/Sujin1135/study-grpc-idl"
)

func (s *server) AddProduct(ctx context.Context, in *pb.Product) {
	// Write your business code
}

func (s *server) GetProduct(ctx context.Context, in *pb.ProductId) {
	// Write your business code
}
