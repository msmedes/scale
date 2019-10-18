package scale

import (
	"fmt"
	"os"

	"github.com/msmedes/scale/internal/pkg/graphql"
	"github.com/msmedes/scale/internal/pkg/node"
	"github.com/msmedes/scale/internal/pkg/rpc"
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
	node := node.NewNode(addr)
	rpcServer := rpc.NewRPC(node)
	graphql := graphql.NewGraphQL(webAddr, rpcServer)

	defer node.Shutdown()

	go graphql.ServerListen()
	go rpcServer.ServerListen()

	if len(join) > 0 {
		node.Join(join)
	}

	go node.StabilizationStart()

	select {}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
