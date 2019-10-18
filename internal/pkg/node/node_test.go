package node

import (
	"testing"

	"github.com/msmedes/scale/internal/pkg/keyspace"
)

const Addr = "0.0.0.0:3000"

func TestNewNode(t *testing.T) {
	n := NewNode(Addr)

	t.Run("sets successor to itself", func(t *testing.T) {
		s, err := n.GetSuccessor()

		if err != nil {
			t.Error(err)
		}

		if !keyspace.Equal(s.GetID(), n.GetID()) {
			t.Errorf("expected successor to be itself, got: %x", n.GetID())
		}
	})

	t.Run("sets predecessor to itself", func(t *testing.T) {
		p, err := n.GetPredecessor()

		if err != nil {
			t.Error(err)
		}

		if !keyspace.Equal(p.GetID(), n.GetID()) {
			t.Errorf("expected predecessor to be itself, got: %x", n.GetID())
		}
	})

	t.Run("sets finger table to itself", func(t *testing.T) {
		ids := n.GetFingerTableIDs()

		for i, id := range ids {
			if !keyspace.Equal(id, n.GetID()) {
				t.Errorf("invalid finger table entry at index %d: %x", i, id)
			}
		}
	})
}

func TestCheckPredecessor(t *testing.T) {
	n := NewNode(Addr)

	t.Run("does nothing when predecessor is itself", func(t *testing.T) {
		n.checkPredecessor()
		p, err := n.GetPredecessor()

		if err != nil {
			t.Error(err)
		}

		if !keyspace.Equal(p.GetID(), n.GetID()) {
			t.Errorf("expected predessor check to not have changed p. got: %x", p.GetID())
		}
	})

	t.Run("removes predecessor if there is an error pinging it", func(t *testing.T) {
		n.predecessor = NewRemoteNode("0.0.0.0:3001", n)
		n.checkPredecessor()
		p, err := n.GetPredecessor()

		if err != nil {
			t.Error(err)
		}

		if p != nil {
			t.Errorf("expected to have removed predecessor")
		}
	})
}
