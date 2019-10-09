package scale

// func (node *Node) fixFingerTable() {
// 	next := 0
// 	ticker := time.NewTicker(10 * time.Second)
// 	for {
// 		select {
// 		case <-ticker.C:
// 			next = node.fixNextFinger(next)
// 			node.logger.Info("fixing finger table")
// 		}
// 	}
// }

// func (node *Node) stabilize() {
// 	ticker := time.NewTicker(10 * time.Second)
// 	for {
// 		select {
// 		case <-ticker.C:
// 			node.logger.Info("stabilizing")
// 		}
// 	}
// }

// func (node *Node) closestPrecedingNode(id Key) *pb.Node {
// 	// I think this could be implemented as binary search?
// 	for i := M - 1; i >= 0; i-- {
// 		finger := node.fingerTable[i]
// 		if finger.RemoteNode == nil {
// 			continue
// 		}

// 		if Between(finger.ID, node.ID, id) {
// 			return finger.RemoteNode
// 		}
// 	}
// 	return node.ProtoNode
// }

// // FindPredecessor find predecessor
// func (node *Node) FindPredecessor(ID Key) Key {
// 	return node.ID
// }

// // GetScaleClient returns a ScaleClient from the node to a specific remote Node.
// func (node *Node) GetScaleClient(remoteNode *RemoteNode) pb.ScaleClient {
// 	// Do we already have a connection
// 	remoteClient, ok := node.clientConns[remoteNode.addr]
// 	if ok {
// 		return remoteClient.client
// 	}
// 	// dial away
// 	conn, err := grpc.Dial(remoteNode.addr, grpc.WithInsecure())

// 	if err != nil {
// 		node.logger.Fatal(err)
// 	}

// 	// Create a new client
// 	client := pb.NewScaleClient(conn)

// 	// Add this client to the map of clientConnections
// 	node.clientConns[remoteNode.Addr] = &clientConnection{client, conn}

// 	return client

// }
