package rpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func serverInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	fmt.Printf("%+v", ctx)

	h, err := handler(ctx, req)

	return h, err
}

func appendTrace(ctx context.Context) context.Context {
	peer, _ := peer.FromContext(ctx)

	ctx = metadata.AppendToOutgoingContext(ctx, "nodes", fmt.Sprintf("%s", peer.Addr))

	return ctx
}

func clientInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	return nil, nil
}
