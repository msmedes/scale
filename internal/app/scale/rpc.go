package scale

import (
	"context"
	"errors"

	pb "github.com/msmedes/scale/internal/app/scale/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// RPC rpc route handler
type RPC struct {
	node   *Node
	logger *zap.SugaredLogger
}

// NewRPC create a new RPC with the given node
func NewRPC(node *Node, logger *zap.SugaredLogger) *RPC {
	return &RPC{
		node:   node,
		logger: logger,
	}
}

// ClosestPrecedingFinger TODO
func (r *RPC) ClosestPrecedingFinger(context.Context, *pb.RemoteQuery) (*pb.RemoteNode, error) {
	return nil, errors.New("not implemented")
}

// FindSuccessor RPC wrapper for node.FindSuccessor
func (r *RPC) FindSuccessor(ctx context.Context, in *pb.RemoteQuery) (*pb.RemoteNode, error) {
	successor := r.node.FindSuccessor(ByteArrayToKey(in.Id))

	res := &pb.RemoteNode{Id: successor.ID[:], Addr: successor.Addr}

	return res, nil
}

// GetSuccessor TODO
func (r *RPC) GetSuccessor(ctx context.Context, in *pb.UpdateReq) (*pb.RemoteNode, error) {
	res := &pb.RemoteNode{
		Id:   r.node.successor.ID[:],
		Addr: r.node.successor.Addr,
	}

	return res, nil
}

// GetPredecessor returns the predecessor of the node
func (r *RPC) GetPredecessor(ctx context.Context, in *pb.UpdateReq) (*pb.RemoteNode, error) {
	predecessor := r.node.GetPredecessor()

	res := &pb.RemoteNode{
		Id:   predecessor.ID[:],
		Addr: predecessor.Addr,
	}

	return res, nil
}

// Notify tells a node that another node (it thinks) it's its predecessor
// man english is a weird language
func (r *RPC) Notify(ctx context.Context, in *pb.RemoteNode) (*pb.RpcOkay, error) {
	message := r.node.Notify(ByteArrayToKey(in.Id), in.Addr)

	// idk is this a thing?
	res := &pb.RpcOkay{
		Message: message,
	}
	r.logger.Info(message)

	return res, nil
}

// SetPredecessor TODO
func (r *RPC) SetPredecessor(context.Context, *pb.UpdateReq) (*pb.RpcOkay, error) {
	return nil, nil
}

// SetSuccessor TODO
func (r *RPC) SetSuccessor(context.Context, *pb.UpdateReq) (*pb.RpcOkay, error) {
	return nil, errors.New("not implemented")
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
