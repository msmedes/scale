package scale

import (
	"context"
	"log"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	pb "github.com/msmedes/scale/internal/app/scale/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// RPC rpc route handler
type RPC struct {
	node   *Node
	logger *zap.Logger
}

// NewRPC create a new RPC with the given node
func NewRPC(node *Node, logger *zap.Logger) *RPC {
	return &RPC{
		node:   node,
		logger: logger,
	}
}

// GetLocal RPC wrapper for node.store.Get
func (r *RPC) GetLocal(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	res := &pb.GetResponse{
		Value: r.node.store.Get(ByteArrayToKey(in.Key)),
	}

	return res, nil
}

// SetLocal RPC wrapper for node.store.Set
func (r *RPC) SetLocal(ctx context.Context, in *pb.SetRequest) (*pb.Success, error) {
	r.node.store.Set(ByteArrayToKey(in.Key), in.Value)

	return &pb.Success{}, nil
}

// FindSuccessor RPC wrapper for node.FindSuccessor
func (r *RPC) FindSuccessor(ctx context.Context, in *pb.RemoteQuery) (*pb.RemoteNode, error) {
	successor := r.node.FindSuccessor(ByteArrayToKey(in.Id))
	res := &pb.RemoteNode{Id: successor.ID[:], Addr: successor.Addr}

	return res, nil
}

// GetSuccessor successor of the node
func (r *RPC) GetSuccessor(context.Context, *pb.Empty) (*pb.RemoteNode, error) {
	successor, err := r.node.GetSuccessor()

	if err != nil {
		return nil, err
	}

	res := &pb.RemoteNode{
		Id:   successor.ID[:],
		Addr: successor.Addr,
	}

	return res, nil
}

// GetPredecessor returns the predecessor of the node
func (r *RPC) GetPredecessor(context.Context, *pb.Empty) (*pb.RemoteNode, error) {
	predecessor, err := r.node.GetPredecessor()

	if err != nil {
		return nil, err
	} else if predecessor == nil {
		empty := &pb.RemoteNode{Present: false}

		return empty, nil
	}

	res := &pb.RemoteNode{
		Id:      predecessor.ID[:],
		Addr:    predecessor.Addr,
		Present: true,
	}

	return res, nil
}

// Notify tells a node that another node (it thinks) it's its predecessor
// man english is a weird language
func (r *RPC) Notify(ctx context.Context, in *pb.RemoteNode) (*pb.Success, error) {
	err := r.node.Notify(ByteArrayToKey(in.Id), in.Addr)

	if err != nil {
		return nil, err
	}

	return &pb.Success{}, nil
}

// GetScaleClient returns a ScaleClient from the node to a specific remote Node.
// I don't really know where to put this.
func GetScaleClient(addr string, node *Node) pb.ScaleClient {
	// Do we already have a connection
	// dial away
	conn, err := grpc.Dial(addr, grpc.WithInsecure())

	if err != nil {
		node.logger.Fatal(err)
	}

	// Create a new client
	client := pb.NewScaleClient(conn)

	if err != nil {
		node.logger.Fatal(err)
	}

	return client
}

// ServerListen start up the server
func (r *RPC) ServerListen() {
	server, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(r.logger),
		)),
	)

	pb.RegisterScaleServer(grpcServer, r)
	grpcServer.Serve(server)
}
