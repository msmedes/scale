package trace

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	pb "github.com/msmedes/scale/internal/pkg/trace/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type traceEntry struct {
	addr         string
	dispatchedAt time.Time
	duration     time.Duration
	functionCall string
	latency      time.Duration
}

// Trace stores the loggers and the map used to store traces
type Trace struct {
	addr   string
	logger *zap.Logger
	port   string
	store  *store
	sugar  *zap.SugaredLogger
}

type store struct {
	store map[string][]*traceEntry
	mutex sync.RWMutex
}

// NewTrace creates a new trace
func NewTrace(addr string, port string) *Trace {
	trace := &Trace{
		addr: addr,
		port: port,
		store: &store{
			store: make(map[string][]*traceEntry),
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

// StartTrace starts a trace, duh
func (t *Trace) StartTrace(ctx context.Context, in *pb.StartTraceRequest) (*pb.Success, error) {
	t.store.mutex.Lock()
	defer t.store.mutex.Unlock()

	traceID := in.TraceID

	if _, ok := t.store.store[traceID]; ok {
		return nil, errors.New("that traceID is already being used...weird")
	}

	dispatchedAt := time.Unix(0, in.Timestamp) //UTC time
	latency := time.Since(dispatchedAt)

	firstEntry := &traceEntry{
		addr:         in.Addr,
		dispatchedAt: dispatchedAt,
		latency:      latency,
		functionCall: in.FunctionCall,
	}

	t.store.store[traceID] = []*traceEntry{firstEntry}

	return &pb.Success{}, nil
}

// AppendTrace appends a node addr to an existing traceID
func (t *Trace) AppendTrace(ctx context.Context, in *pb.AppendTraceRequest) (*pb.Success, error) {
	t.store.mutex.Lock()
	defer t.store.mutex.Unlock()

	traceID := in.TraceID
	addr := in.Addr
	if in.FunctionCall == "SetLocal" {
		fmt.Println("lol what")
	}

	if _, ok := t.store.store[traceID]; !ok {
		return nil, errors.New("that traceID does not exist in the store")
	}

	trace := t.store.store[traceID]
	lastEntry := trace[len(trace)-1]

	if lastEntry.addr != addr {
		dispatchedAt := time.Unix(0, in.Timestamp)
		latency := time.Since(dispatchedAt)

		lastEntry.duration = dispatchedAt.Sub(lastEntry.dispatchedAt)
		entry := &traceEntry{
			addr:         addr,
			dispatchedAt: dispatchedAt,
			latency:      latency,
			functionCall: in.FunctionCall,
		}
		t.store.store[traceID] = append(t.store.store[traceID], entry)
	}

	return &pb.Success{}, nil
}

// GetTrace returns the trace from a given TraceID and then
// deletes the trace from the store
func (t *Trace) GetTrace(ctx context.Context, in *pb.TraceQuery) (*pb.TraceMessage, error) {
	t.store.mutex.RLock()
	defer t.store.mutex.RUnlock()

	traceID := in.TraceID

	trace, ok := t.store.store[traceID]
	if !ok {
		return nil, errors.New("that traceID does not exist in the store")
	}
	dispatchedAt := time.Unix(0, in.Timestamp)
	lastEntry := trace[len(trace)-1]

	lastEntry.duration = dispatchedAt.Sub(lastEntry.dispatchedAt)
	for _, entry := range trace {
		fmt.Printf("%+v\n", entry)
	}
	traceMessage := &pb.TraceMessage{
		Trace: createTraceEntries(trace),
	}
	delete(t.store.store, traceID)

	return traceMessage, nil
}

func createTraceEntries(trace []*traceEntry) []*pb.TraceEntry {
	var entries []*pb.TraceEntry

	for _, entry := range trace {
		pbEntry := &pb.TraceEntry{
			Addr:         entry.addr,
			FunctionCall: entry.functionCall,
			Duration:     entry.duration.String(),
		}
		entries = append(entries, pbEntry)
	}
	return entries
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
