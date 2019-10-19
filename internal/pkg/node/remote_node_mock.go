package node

import (
	"github.com/msmedes/scale/internal/pkg/scale"
)

// RemoteNodeMock mock remote node
type RemoteNodeMock struct {
	scale.RemoteNode

	ID                      scale.Key
	Addr                    string
	findPredecessorResponse scale.RemoteNode
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
func (m *RemoteNodeMock) FindPredecessor(key scale.Key) (scale.RemoteNode, error) {
	return m.findPredecessorResponse, nil
}

//Notify notify
func (m *RemoteNodeMock) Notify(node scale.Node) error {
	return nil
}
