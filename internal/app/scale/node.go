package scale

import (
	"context"
	"time"

	uuid "github.com/google/uuid"
	pb "github.com/msmedes/scale/internal/app/scale/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Node main node class
type Node struct {
	ID          Key
	Addr        string
	predecessor *Node
	successor   *Node
	fingerTable FingerTable
	store       *Store
	logger      *zap.SugaredLogger
}

// RemoteNode contains metadata about another node
type RemoteNode struct {
	ID   Key
	Addr string
}

// NewNode create a new node
func NewNode(addr string, logger *zap.SugaredLogger) *Node {
	node := &Node{
		ID:     genID(),
		Addr:   addr,
		store:  NewStore(),
		logger: logger,
	}

	node.fingerTable = NewFingerTable(M, node)
	node.successor = node

	return node
}

// Join join an existing network via another node
func (node *Node) Join(addr string) {
	node.logger.Infof("joining network via node at %s", addr)

	conn, err := grpc.Dial(addr, grpc.WithInsecure())

	if err != nil {
		node.logger.Fatal(err)
	}

	defer conn.Close()

	client := pb.NewScaleClient(conn)

	successor, err := client.FindSuccessor(
		context.Background(),
		&pb.RemoteQuery{Id: KeyToString(node.ID)},
	)

	if err != nil {
		node.logger.Fatal(err)
	}

	node.logger.Infof("found successor: %s", successor.Id)
	node.logger.Info("joined network")
}

// Shutdown leave the network
func (node *Node) Shutdown() {
	node.logger.Info("exiting")
}

func (node *Node) stabilize(ticker *time.Ticker) {
	node.logger.Fatal("not implemented")
}

func (node *Node) notify(remoteNode *RemoteNode) {
	node.logger.Fatal("not implemented")
}

func (node *Node) findSuccessor(ID []byte) {
	node.logger.Fatal("not implemented")
}

func (node *Node) findPredecessor(ID []byte) {
	node.logger.Fatal("not implemented")
}

func genID() Key {
	ID, _ := uuid.NewRandom()
	return GenerateKey(ID.String())
}
