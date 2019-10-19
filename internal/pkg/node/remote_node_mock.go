package node

import (
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
func (m *RemoteNodeMock) FindPredecessor(scale.Key) (scale.RemoteNode, error) {
	return m.findPredecessorResponse, nil
}

//FindSuccessor finds the predecessor to the id
func (m *RemoteNodeMock) FindSuccessor(scale.Key) (scale.RemoteNode, error) {
	return m.findSuccessorResponse, nil
}

//GetPredecessor finds the predecessor to the id
func (m *RemoteNodeMock) GetPredecessor() (scale.RemoteNode, error) {
	return m.getPredecessorResponse, nil
}

//GetSuccessor finds the predecessor to the id
func (m *RemoteNodeMock) GetSuccessor() (scale.RemoteNode, error) {
	return m.getSuccessorResponse, nil
}

//Notify notify
func (m *RemoteNodeMock) Notify(scale.Node) error {
	return nil
}

//ClosestPrecedingFinger finds the predecessor to the id
func (m *RemoteNodeMock) ClosestPrecedingFinger(scale.Key) (scale.RemoteNode, error) {
	return m.closestPrecedingFingerResponse, nil
}
