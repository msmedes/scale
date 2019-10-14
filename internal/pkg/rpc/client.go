package rpc

import (
	"log"

	pb "github.com/msmedes/scale/internal/pkg/rpc/proto"
	"google.golang.org/grpc"
)

// NewClient returns a ScaleClient from the node to a specific remote Node.
func NewClient(addr string) (pb.ScaleClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())

	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewScaleClient(conn)

	if err != nil {
		log.Fatal(err)
	}

	return client, conn
}
