package scale

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
)

var (
	port    = getEnv("PORT", "3000")
	addr    = fmt.Sprintf("0.0.0.0:%s", port)
	join    = getEnv("JOIN", "")
	webPort = getEnv("WEB", "8000")
	webAddr = fmt.Sprintf("0.0.0.0:%s", webPort)
)

// ServerListen create Node to represent this current node. Start up a grpc
// server to accept requests from remote nodes and invoke methods on the node object
// also, fire up the background process to periodically stabilize the node's finger table
func ServerListen() {
	logger, err := zap.NewDevelopment()

	defer logger.Sync()

	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	sugar := logger.Sugar()
	node := NewNode(addr, sugar)
	rpc := NewRPC(node, logger)
	graphql := NewGraphQL(webAddr, sugar, rpc)

	defer node.Shutdown()

	go graphql.ServerListen()
	sugar.Infof("listening - graphql: %s", webAddr)

	go rpc.ServerListen()
	sugar.Infof("listening - internode: %s", addr)

	go node.StabilizationStart()

	if len(join) > 0 {
		node.Join(join)
	}

	sugar.Infof("node.id: %s", KeyToString(node.ID))

	for {
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
