package scale

import (
	"context"
	"errors"
	"log"

	pb "github.com/msmedes/scale/internal/app/scale/proto"
)

type RPC struct {
	node *Node
}

func NewRPC(node *Node) *RPC {
	return &RPC{
		node: node,
	}
}

func (r *RPC) ClosestPrecedingFinger(context.Context, *pb.RemoteQuery) (*pb.IdReply, error) {
	return nil, errors.New("not implemented")
}

func (r *RPC) FindSuccessor(context.Context, *pb.RemoteQuery) (*pb.IdReply, error) {
	return nil, errors.New("not implemented")
}

func (r *RPC) GetSuccessorId(ctx context.Context, in *pb.RemoteId) (*pb.IdReply, error) {
	log.Printf("%s", in.Id)

	res := &pb.IdReply{
		Id: IdToString(r.node.Id),
	}

	return res, nil
}

func (r *RPC) Notify(context.Context, *pb.RemoteNode) (*pb.RpcOkay, error) {
	return nil, errors.New("not implemented")
}

func (r *RPC) SetPredecessorId(context.Context, *pb.UpdateReq) (*pb.RpcOkay, error) {
	return nil, nil
}

func (r *RPC) SetSuccessorId(context.Context, *pb.UpdateReq) (*pb.RpcOkay, error) {
	return nil, errors.New("not implemented")
}
