package scale

import (
	"context"
	"errors"

	pb "github.com/msmedes/scale/internal/app/scale/proto"
)

type RPC struct{}

func (r *RPC) ClosestPrecedingFinger(context.Context, *pb.RemoteQuery) (*pb.IdReply, error) {
	return nil, errors.New("not implemented")
}

func (r *RPC) FindSuccessor(context.Context, *pb.RemoteQuery) (*pb.IdReply, error) {
	return nil, errors.New("not implemented")
}

func (r *RPC) GetSuccessorId(context.Context, *pb.RemoteId) (*pb.IdReply, error) {
	return nil, errors.New("not implemented")
}

func (r *RPC) Notify(context.Context, *pb.RemoteNode) (*pb.RpcOkay, error) {
	return nil, errors.New("not implemented")
}

func (r *RPC) SetPredecessorId(context.Context, *pb.UpdateReq) (*pb.RpcOkay, error) {
	return nil, errors.New("not implemented")
}

func (r *RPC) SetSuccessorId(context.Context, *pb.UpdateReq) (*pb.RpcOkay, error) {
	return nil, errors.New("not implemented")
}
