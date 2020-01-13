package rpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// TraceServerInterceptor is the interceptor for the server trace
func TraceServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)

		ctx = context.WithValue(ctx, "trace", md["trace"])
		fmt.Printf("***SERVER CONTEXT***\n %+v\n", ctx)
		return handler(ctx, req)
	}
}

// TraceClientInterceptor doesn't do much yet
func TraceClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error {

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
