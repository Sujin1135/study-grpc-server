package main

import (
	"google.golang.org/grpc"
	"log"
	"time"
)

type wrappedStream struct {
	grpc.ServerStream
}

func PrintLog(m interface{}, msg string) {
	log.Printf("===== [Server Interceptor Wrapper] "+
		msg+" (Type: %T) at %s",
		m, time.Now().Format(time.RFC3339))
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	PrintLog(m, "Receive a message")
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	PrintLog(m, "Send a message")
	return w.ServerStream.SendMsg(m)
}

func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

func orderServerStreamInterceptor(srv interface{},
	ss grpc.ServerStream, info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	log.Printf("===== [Server Stream Interceptor] ", info.FullMethod)

	err := handler(srv, newWrappedStream(ss))

	if err != nil {
		log.Printf("RPC failed with error %v", err)
	}
	return err
}
