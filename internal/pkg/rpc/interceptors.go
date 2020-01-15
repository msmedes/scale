package rpc

import (
	"context"

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
		// fmt.Printf("SERVER MD: %s: %+v\n", info.FullMethod, ctx)

		ctx = metadata.NewIncomingContext(ctx, md)

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

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.Pairs()
		}

		ctx = metadata.NewOutgoingContext(ctx, md)
		// fmt.Printf("CLIENT CONTEXT %s %+v\n", method, ctx)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
