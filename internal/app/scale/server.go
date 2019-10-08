package scale

import (
	"fmt"
	"log"
	"net"
	"os"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/msmedes/scale/internal/app/scale/proto"
)

var (
	port = getEnv("PORT", "3000")
	addr = fmt.Sprintf("0.0.0.0:%s", port)
	join = getEnv("JOIN", "")
)

func ServerListen() {
	logger, err := zap.NewDevelopment()

	defer logger.Sync()

	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	sugar := logger.Sugar()
	node := NewNode(addr, sugar)

	defer node.Shutdown()

	if len(join) > 0 {
		node.Join(join)
	}

	sugar.Infof("listening: %s", addr)
	sugar.Infof("node.id: %s", KeyToString(node.ID))
	// log.Printf("node.fingerTable: %s", node.fingerTable)

	startGRPC(node, logger)
}

func startGRPC(node *Node, logger *zap.Logger) {
	server, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	rpc := NewRPC(node)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(logger),
		)),
	)

	pb.RegisterScaleServer(grpcServer, rpc)
	grpcServer.Serve(server)
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
