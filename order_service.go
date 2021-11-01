package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"strings"
	pb "study-grpc-server/order/order"
)

const (
	orderBatchSize = 3
)

var orderMap = make(map[string]pb.Order)

type orderServer struct {
	orderMap map[string]*pb.Order
}

func (s *orderServer) AddOrder(ctx context.Context, order *pb.Order) (*wrappers.StringValue, error) {
	orderMap[order.Id] = *order
	log.Printf("Order %v : %v - Added.", order.Id, order.Description)
	return &wrappers.StringValue{Value: order.Id}, nil
}

func (s *orderServer) GetOrders(searchQuery *wrappers.StringValue, stream pb.OrderManagement_GetOrdersServer) error {
	for key, order := range orderMap {
		log.Print(key, order)
		for _, itemStr := range order.Items {
			log.Print(itemStr)
			if strings.Contains(itemStr, searchQuery.Value) {
				err := stream.Send(&order)
				if err != nil {
					return fmt.Errorf("error sending message to stream: %v", err)
				}
				log.Print("Matching order Found: " + key)
				break
			}
		}
	}
	return nil
}

func (s *orderServer) UpdateOrders(stream pb.OrderManagement_UpdateOrdersServer) error {
	ordersStr := "Updated Order IDs : "
	for {
		order, err := stream.Recv()
		if err == io.EOF {
			// Finished reading the order stream.
			return stream.SendAndClose(&wrappers.StringValue{Value: "Orders processed " + ordersStr})
		}
		// Update order
		orderMap[order.Id] = *order

		log.Printf("Order ID %v : Updated", order.Id)
		ordersStr += order.Id + ", "

		md, ok := metadata.FromIncomingContext(stream.Context())

		log.Printf("stream gRPC md is %v", md)
		log.Printf("stream gRPC ok is %v", ok)
	}
}

func initSampleData() {
	orderMap["102"] = pb.Order{Id: "102", Items: []string{"Google Pixel 3A", "Mac Book Pro"}, Destination: "Mountain View, CA", Price: 1800.00}
	orderMap["103"] = pb.Order{Id: "103", Items: []string{"Apple Watch S4"}, Destination: "San Jose, CA", Price: 400.00}
	orderMap["104"] = pb.Order{Id: "104", Items: []string{"Google Home Mini", "Google Nest Hub"}, Destination: "Mountain View, CA", Price: 400.00}
	orderMap["105"] = pb.Order{Id: "105", Items: []string{"Amazon Echo"}, Destination: "San Jose, CA", Price: 30.00}
	orderMap["106"] = pb.Order{Id: "106", Items: []string{"Amazon Echo", "Apple iPhone XS"}, Destination: "Mountain View, CA", Price: 300.00}
}

func (s *orderServer) ProcessOrders(stream pb.OrderManagement_ProcessOrdersServer) error {
	batchMarker := 1
	combinedShipmentMap := make(map[string]pb.CombinedShipment)

	for {
		orderId, err := stream.Recv()
		log.Printf("Reading Proc order: %s", orderId)

		if err == io.EOF {
			log.Printf("EOF : %s", orderId)
			for _, comb := range combinedShipmentMap {
				log.Fatalf("comb: %v", comb)
				if err := stream.Send(&comb); err != nil {
					return err
				}
			}
			return nil
		}
		if err != nil {
			return err
		}

		// business logics as below
		destination := orderMap[orderId.GetValue()].Destination
		shipment, found := combinedShipmentMap[destination]

		if found {
			ord := orderMap[orderId.GetValue()]
			shipment.OrderList = append(shipment.OrderList, &ord)
			combinedShipmentMap[destination] = shipment
		} else {
			comShip := pb.CombinedShipment{Id: "cmb - " + (orderMap[orderId.GetValue()].Destination), Status: "Processed!"}
			ord := orderMap[orderId.GetValue()]
			comShip.OrderList = append(shipment.OrderList, &ord)
			combinedShipmentMap[destination] = comShip
			log.Print(len(comShip.OrderList), comShip.GetId())
		}

		// end business logics

		if batchMarker == orderBatchSize {
			for _, comb := range combinedShipmentMap {
				stream.Send(&comb)
			}
			batchMarker = 0
			combinedShipmentMap = make(map[string]pb.CombinedShipment)
		} else {
			batchMarker++
		}
	}
}
