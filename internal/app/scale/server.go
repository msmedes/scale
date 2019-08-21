package scale

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/msmedes/scale/internal/app/scale/proto"
)

func ServerListen() {
	port := flag.Int("port", 3000, "port for grpc server")
	flag.Parse()

	server, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("listening on: 0.0.0.0:%d", *port)
	grpcServer := grpc.NewServer()
	pb.RegisterScaleServer(grpcServer, &RPC{})
	grpcServer.Serve(server)
}
