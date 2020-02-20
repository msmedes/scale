package node

import (
	"context"

	"github.com/msmedes/scale/internal/pkg/scale"
)

// RemoteNodeMock mock remote node
type RemoteNodeMock struct {
	scale.RemoteNode

	ID                             scale.Key
	Addr                           string
	findPredecessorResponse        scale.RemoteNode
	findSuccessorResponse          scale.RemoteNode
	getSuccessorResponse           scale.RemoteNode
	getPredecessorResponse         scale.RemoteNode
	closestPrecedingFingerResponse scale.RemoteNode
}

// GetID getter for ID
func (m *RemoteNodeMock) GetID() scale.Key {
	return m.ID
}

// GetAddr getter for address
func (m *RemoteNodeMock) GetAddr() string {
	return m.Addr
}

//FindPredecessor finds the predecessor to the id
func (m *RemoteNodeMock) FindPredecessor(context.Context, scale.Key) (scale.RemoteNode, error) {
	return m.findPredecessorResponse, nil
}

//FindSuccessor finds the predecessor to the id
func (m *RemoteNodeMock) FindSuccessor(context.Context, scale.Key) (scale.RemoteNode, error) {
	return m.findSuccessorResponse, nil
}

//GetPredecessor finds the predecessor to the id
func (m *RemoteNodeMock) GetPredecessor(context.Context) (scale.RemoteNode, error) {
	return m.getPredecessorResponse, nil
}

//GetSuccessor finds the predecessor to the id
func (m *RemoteNodeMock) GetSuccessor(context.Context) (scale.RemoteNode, error) {
	return m.getSuccessorResponse, nil
}

//Notify notify
func (m *RemoteNodeMock) Notify(scale.Node) error {
	return nil
}

//ClosestPrecedingFinger finds the predecessor to the id
func (m *RemoteNodeMock) ClosestPrecedingFinger(context.Context, scale.Key) (scale.RemoteNode, error) {
	return m.closestPrecedingFingerResponse, nil
}
