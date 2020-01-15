package trace

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	pb "github.com/msmedes/scale/internal/pkg/trace/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Trace stores the loggers and the map used to store traces
type Trace struct {
	addr   string
	logger *zap.Logger
	port   string
	sugar  *zap.SugaredLogger
	store  *store
}

type store struct {
	store map[string][]string
	mutex sync.RWMutex
}

// NewTrace creates a new trace
func NewTrace(addr string, port string) *Trace {
	trace := &Trace{
		addr: addr,
		port: port,
		store: &store{
			store: make(map[string][]string),
		},
	}

	logger, err := zap.NewDevelopment(
		zap.Fields(
			zap.String("node", fmt.Sprintf("%s", addr)),
		),
	)

	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	sugar := logger.Sugar()
	trace.logger = logger
	trace.sugar = sugar

	return trace
}

// GetAddr returns the address for the trace server
func (t *Trace) GetAddr() string {
	return t.addr
}

// ServerListen starts up the trace server
func (t *Trace) ServerListen() {

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(t.logger),
		)),
	}

	server, err := net.Listen("tcp", t.GetAddr())

	if err != nil {
		t.sugar.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterTraceServer(grpcServer, t)

	t.sugar.Infof("listening: %s", t.GetAddr())

	grpcServer.Serve(server)
}

// Shutdown the trace server
func (t *Trace) Shutdown() {
	t.logger.Sync()
	t.sugar.Sync()
}

// StartTrace starts a trace, duh
func (t *Trace) StartTrace(ctx context.Context, in *pb.StartTraceRequest) (*pb.Success, error) {
	t.store.mutex.Lock()
	defer t.store.mutex.Unlock()

	traceID := in.TraceID

	if _, ok := t.store.store[traceID]; ok {
		return nil, errors.New("that traceID is already being used...weird")
	}

	t.store.store[traceID] = []string{}

	t.sugar.Infof("trace %+v", t.store.store)

	return &pb.Success{}, nil
}

// AppendTrace appends a node addr to an existing traceID
func (t *Trace) AppendTrace(ctx context.Context, in *pb.AppendTraceRequest) (*pb.Success, error) {
	t.store.mutex.Lock()
	defer t.store.mutex.Unlock()

	t.sugar.Infof("PRE-APPEND %+v\n", t.store.store)

	traceID := in.TraceID
	node := in.Addr

	if _, ok := t.store.store[traceID]; !ok {
		return nil, errors.New("that traceID does not exist in the store")
	}

	trace := t.store.store[traceID]
	if len(trace) == 0 {
		t.store.store[traceID] = append(t.store.store[traceID], node)
	} else {
		if trace[len(trace)-1] != node {
			t.store.store[traceID] = append(t.store.store[traceID], node)
		}
	}

	t.sugar.Infof("POST APPEND %+v\n", t.store.store)

	return &pb.Success{}, nil
}

// GetTrace returns the trace from a given TraceID.
// May extend to delete the trace info, or perhaps store maybe 20
// traces in memory.
func (t *Trace) GetTrace(ctx context.Context, in *pb.TraceQuery) (*pb.TraceMessage, error) {
	t.store.mutex.RLock()
	defer t.store.mutex.RUnlock()

	traceID := in.TraceID

	trace, ok := t.store.store[traceID]
	if !ok {
		return nil, errors.New("that traceID does not exist in the store")
	}
	t.sugar.Infof("%+v", t.store.store)

	traceMessage := &pb.TraceMessage{
		Trace: trace,
	}
	delete(t.store.store, traceID)

	return traceMessage, nil
}
