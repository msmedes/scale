package rpc

import (
	"log"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"google.golang.org/grpc/keepalive"

	pb "github.com/msmedes/scale/internal/pkg/rpc/proto"
	"google.golang.org/grpc"
)

// NewClient returns a ScaleClient from the node to a specific remote Node.
func NewClient(addr string) (pb.ScaleClient, *grpc.ClientConn) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(TraceClientInterceptor())),
	}

	conn, err := grpc.Dial(addr, opts...)

	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewScaleClient(conn)

	return client, conn
}
