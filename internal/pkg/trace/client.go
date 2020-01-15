package trace

import (
	"log"

	"google.golang.org/grpc"

	pb "github.com/msmedes/scale/internal/pkg/trace/proto"
)

// NewClient returns a new TraceClient the specified trace Node
func NewClient(addr string) (pb.TraceClient, *grpc.ClientConn) {

	conn, err := grpc.Dial(addr, grpc.WithInsecure())

	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewTraceClient(conn)

	return client, conn
}
