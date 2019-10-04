package scale

import (
	"context"
	"time"

	uuid "github.com/google/uuid"
	pb "github.com/msmedes/scale/internal/app/scale/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Node struct {
	Id          Key
	Addr        string
	predecessor *Node
	successor   *Node
	fingerTable FingerTable
	store       *Store
	logger      *zap.SugaredLogger
}

type RemoteNode struct {
	Id   Key
	Addr string
}

func NewNode(addr string, logger *zap.SugaredLogger) *Node {
	node := &Node{
		Id:     genId(),
		Addr:   addr,
		store:  NewStore(),
		logger: logger,
	}

	node.fingerTable = NewFingerTable(M, node)
	node.successor = node

	return node
}

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
		&pb.RemoteQuery{Id: KeyToString(node.Id)},
	)

	if err != nil {
		node.logger.Fatal(err)
	}

	node.logger.Infof("found successor: %s", successor.Id)
	node.logger.Info("joined network")
}

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

func genId() Key {
	Id, _ := uuid.NewRandom()
	return GenerateKey(Id.String())
}
