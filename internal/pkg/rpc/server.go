package rpc

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc/keepalive"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/msmedes/scale/internal/pkg/keyspace"
	pb "github.com/msmedes/scale/internal/pkg/rpc/proto"
	"github.com/msmedes/scale/internal/pkg/scale"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// RPC rpc route handler
type RPC struct {
	node   scale.Node
	sugar  *zap.SugaredLogger
	logger *zap.Logger
}

// NewRPC create a new rpc
func NewRPC(node scale.Node) *RPC {
	rpc := &RPC{node: node}

	logger, err := zap.NewDevelopment(
		zap.Fields(
			zap.String("node", keyspace.KeyToString(node.GetID())),
		),
	)

	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	sugar := logger.Sugar()
	rpc.logger = logger
	rpc.sugar = sugar

	return rpc
}

// TransferKeys proxy to node.TransferKeys
func (r *RPC) TransferKeys(ctx context.Context, in *pb.KeyTransferRequest) (*pb.Success, error) {
	r.node.TransferKeys(keyspace.ByteArrayToKey(in.GetId()), in.GetAddr())
	return &pb.Success{}, nil
}

// Ping health check
func (r *RPC) Ping(ctx context.Context, in *pb.Empty) (*pb.Success, error) {
	return &pb.Success{}, nil
}

// Get rpc wrapper for node.Get
func (r *RPC) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	val, err := r.node.Get(keyspace.ByteArrayToKey(in.GetKey()))

	if err != nil {
		return nil, err
	}

	res := &pb.GetResponse{Value: val}

	return res, nil
}

// SetLocal rpc wrapper for node.store.Set
func (r *RPC) SetLocal(ctx context.Context, in *pb.SetRequest) (*pb.Success, error) {
	r.node.SetLocal(keyspace.ByteArrayToKey(in.GetKey()), in.GetValue())

	return &pb.Success{}, nil
}

// FindSuccessor rpc wrapper for node.FindSuccessor
func (r *RPC) FindSuccessor(ctx context.Context, in *pb.RemoteQuery) (*pb.RemoteNode, error) {
	successor, err := r.node.FindSuccessor(keyspace.ByteArrayToKey(in.Id))

	if err != nil {
		return nil, err
	}

	id := successor.GetID()

	res := &pb.RemoteNode{
		Id:   id[:],
		Addr: successor.GetAddr(),
	}

	return res, nil
}

// ClosestPrecedingFinger returns the node that is the closest predecessor
// of the ID
func (r *RPC) ClosestPrecedingFinger(ctx context.Context, in *pb.RemoteQuery) (*pb.RemoteNode, error) {
	closestPrecedingFinger, err := r.node.ClosestPrecedingFinger(keyspace.ByteArrayToKey(in.Id))

	if err != nil {
		return nil, err
	}

	id := closestPrecedingFinger.GetID()

	res := &pb.RemoteNode{
		Id:   id[:],
		Addr: closestPrecedingFinger.GetAddr(),
	}

	return res, nil
}

// FindPredecessor returns the predecessor of the input key
func (r *RPC) FindPredecessor(ctx context.Context, in *pb.RemoteQuery) (*pb.RemoteNode, error) {
	predecessor, err := r.node.FindPredecessor(keyspace.ByteArrayToKey(in.Id))

	if err != nil {
		return nil, err
	}

	id := predecessor.GetID()

	res := &pb.RemoteNode{
		Id:   id[:],
		Addr: predecessor.GetAddr(),
	}

	return res, nil
}

// GetSuccessor successor of the node
func (r *RPC) GetSuccessor(context.Context, *pb.Empty) (*pb.RemoteNode, error) {
	successor, err := r.node.GetSuccessor()

	if err != nil {
		return nil, err
	}

	id := successor.GetID()

	res := &pb.RemoteNode{
		Id:   id[:],
		Addr: successor.GetAddr(),
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

	id := predecessor.GetID()

	res := &pb.RemoteNode{
		Id:      id[:],
		Addr:    predecessor.GetAddr(),
		Present: true,
	}

	return res, nil
}

// Notify tells a node that another node (it thinks) it's its predecessor
// man english is a weird language
func (r *RPC) Notify(ctx context.Context, in *pb.RemoteNode) (*pb.Success, error) {
	err := r.node.Notify(keyspace.ByteArrayToKey(in.Id), in.Addr)

	if err != nil {
		return nil, err
	}

	return &pb.Success{}, nil
}

// Set rpc wrapper for node.Set
func (r *RPC) Set(ctx context.Context, in *pb.SetRequest) (*pb.Success, error) {
	r.node.Set(keyspace.ByteArrayToKey(in.GetKey()), in.GetValue())

	return &pb.Success{}, nil
}

// GetNodeMetadata return metadata about this node
func (r *RPC) GetNodeMetadata(context.Context, *pb.Empty) (*pb.NodeMetadata, error) {
	id := r.node.GetID()

	var ft [][]byte

	for _, k := range r.node.GetFingerTableIDs() {
		ft = append(ft, k[:])
	}

	meta := &pb.NodeMetadata{
		Id:          id[:],
		Addr:        r.node.GetAddr(),
		Port:        r.node.GetPort(),
		FingerTable: ft,
	}

	predecessor, err := r.node.GetPredecessor()

	if err != nil {
		return nil, err
	}

	if predecessor != nil {
		predID := predecessor.GetID()
		meta.PredecessorId = predID[:]
		meta.PredecessorAddr = predecessor.GetAddr()
	}

	successor, err := r.node.GetSuccessor()

	if err != nil {
		return nil, err
	}

	if successor != nil {
		succID := successor.GetID()
		meta.SuccessorId = succID[:]
		meta.SuccessorAddr = successor.GetAddr()
	}

	return meta, nil
}

// GetLocal rpc wrapper for node.store.Get
func (r *RPC) GetLocal(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	val, err := r.node.GetLocal(keyspace.ByteArrayToKey(in.GetKey()))

	if err != nil {
		return nil, err
	}

	res := &pb.GetResponse{Value: val}

	return res, nil
}

// ServerListen start up the server
func (r *RPC) ServerListen() {
	opts := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(r.logger),
		)),
	}

	server, err := net.Listen("tcp", r.node.GetAddr())

	if err != nil {
		r.sugar.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterScaleServer(grpcServer, r)

	r.sugar.Infof("listening: %s", r.node.GetAddr())

	grpcServer.Serve(server)
}

//Shutdown clean exit
func (r *RPC) Shutdown() {
	r.logger.Sync()
	r.sugar.Sync()
}
