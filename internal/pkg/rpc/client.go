package rpc

import (
	"log"
	"time"

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
	}

	conn, err := grpc.Dial(addr, opts...)

	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewScaleClient(conn)

	if err != nil {
		log.Fatal(err)
	}

	return client, conn
}
