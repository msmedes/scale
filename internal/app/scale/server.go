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
	node := NewNode(addr)
	logger, err := zap.NewDevelopment(zap.Fields(zap.String("node", KeyToString(node.ID))))

	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	sugar := logger.Sugar()
	node.logger = sugar
	rpc := NewRPC(node, logger)
	graphql := NewGraphQL(webAddr, sugar, rpc)

	if len(join) > 0 {
		node.Join(join)
	}

	defer node.Shutdown()
	defer logger.Sync()

	go node.StabilizationStart()
	go graphql.ServerListen()
	go rpc.ServerListen()

	sugar.Infof("listening - graphql: %s", webAddr)
	sugar.Infof("listening - internode: %s", addr)

	select {}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
