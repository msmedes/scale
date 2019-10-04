package scale

import (
	"context"
	"errors"

	pb "github.com/msmedes/scale/internal/app/scale/proto"
)

// RPC rpc route handler
type RPC struct {
	node *Node
}

// NewRPC create a new RPC with the given node
func NewRPC(node *Node) *RPC {
	return &RPC{
		node: node,
	}
}

// ClosestPrecedingFinger TODO
func (r *RPC) ClosestPrecedingFinger(context.Context, *pb.RemoteQuery) (*pb.IdReply, error) {
	return nil, errors.New("not implemented")
}

// FindSuccessor TODO
func (r *RPC) FindSuccessor(ctx context.Context, in *pb.RemoteQuery) (*pb.IdReply, error) {
	res := &pb.IdReply{
		Id: KeyToString(r.node.ID),
	}

	return res, nil
}

// GetSuccessor TODO
func (r *RPC) GetSuccessor(ctx context.Context, in *pb.RemoteId) (*pb.IdReply, error) {
	res := &pb.IdReply{
		Id: KeyToString(r.node.ID),
	}

	return res, nil
}

// Notify TODO
func (r *RPC) Notify(context.Context, *pb.RemoteNode) (*pb.RpcOkay, error) {
	return nil, errors.New("not implemented")
}

// SetPredecessor TODO
func (r *RPC) SetPredecessor(context.Context, *pb.UpdateReq) (*pb.RpcOkay, error) {
	return nil, nil
}

// SetSuccessor TODO
func (r *RPC) SetSuccessor(context.Context, *pb.UpdateReq) (*pb.RpcOkay, error) {
	return nil, errors.New("not implemented")
}
