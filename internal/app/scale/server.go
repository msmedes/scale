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

	server, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	node := NewNode()
	rpc := NewRPC(node)
	pb.RegisterScaleServer(grpcServer, rpc)

	defer grpcServer.Serve(server)

	log.Printf("listening on: :%d", *port)
	log.Printf("node.id: %s", IdToString(node.ID))
	log.Printf("fingerTable: %s", node.fingerTable)
}
