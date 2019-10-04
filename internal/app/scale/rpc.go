package scale

import (
	"context"
	"errors"

	pb "github.com/msmedes/scale/internal/app/scale/proto"
	"go.uber.org/zap"
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

// ClosestPrecedingNode TODO
func (r *RPC) ClosestPrecedingNode(context.Context, *pb.RemoteQuery) (*pb.IdReply, error) {
	return nil, errors.New("not implemented")
}

// FindSuccessor RPC wrapper for node.FindSuccessor
func (r *RPC) FindSuccessor(ctx context.Context, in *pb.RemoteQuery) (*pb.IdReply, error) {
	id := r.node.FindSuccessor(StringToKey(in.Id))
	res := &pb.IdReply{Id: KeyToString(id)}

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
