package node

import (
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/msmedes/scale/internal/pkg/keyspace"
)

const (
	Addr1 = "0.0.0.0:3000"
	Addr2 = "0.0.0.0:3001"
)

func newAddr(t *testing.T) string {
	id, err := uuid.NewRandom()

	if err != nil {
		log.Fatal(err)
	}

	return id.String()
}

func TestNewNode(t *testing.T) {
	clearRemotes()
	n := NewNode(Addr1)

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
	clearRemotes()
	n := NewNode(Addr1)

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
		n.predecessor = newRemoteNode(Addr2)
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

func TestJoin(t *testing.T) {
	clearRemotes()
	t.Run("one other node in network", func(t *testing.T) {
		n2 := NewNode(Addr2)

		n1 := &RemoteNodeMock{
			Addr: Addr1,
			ID:   keyspace.GenerateKey(Addr1),
		}

		n1.findSuccessorResponse = n1
		n1.findPredecessorResponse = n1
		n1.getSuccessorResponse = n1
		n1.getPredecessorResponse = n1

		n2.join(n1)

		t.Run("sets predecessor to other node", func(t *testing.T) {
			if n2.predecessor.GetID() != n1.GetID() {
				t.Errorf("expected n2.predecessor to be n1. got: %x", n2.predecessor.GetID())
			}
		})

		t.Run("sets successor to other node", func(t *testing.T) {
			if n2.successor.GetID() != n1.GetID() {
				t.Errorf("expected n2.successor to be n1. got: %x", n2.successor.GetID())
			}
		})

		t.Run("sets finger table to other node", func(t *testing.T) {
			ids := n2.GetFingerTableIDs()

			for i, id := range ids {
				if !keyspace.Equal(id, n1.GetID()) {
					t.Errorf("invalid finger table entry at index %d: %x", i, id)
				}
			}
		})
	})
}

func TestClosestPrecedingFinger(t *testing.T) {
	clearRemotes()

	n := &Node{addr: newAddr(t), id: [4]byte{1}}
	nRemote := n.toRemoteNode()
	n.fingerTable = newFingerTable(nRemote)

	t.Run("ft is all one value < key", func(t *testing.T) {
		key := [4]byte{2}
		closest, err := n.ClosestPrecedingFinger(key)

		if err != nil {
			t.Error(err)
		}

		if !keyspace.Equal(closest.GetID(), n.GetID()) {
			t.Errorf("expected %x got %x", n.GetID(), closest.GetID())
		}
	})

	t.Run("ft is all one value > key", func(t *testing.T) {
		key := [4]byte{0}
		closest, err := n.ClosestPrecedingFinger(key)

		if err != nil {
			t.Error(err)
		}

		if !keyspace.Equal(closest.GetID(), n.GetID()) {
			t.Errorf("expected %x got %x", n.GetID(), closest.GetID())
		}
	})

	t.Run("ft is all one value == key", func(t *testing.T) {
		key := [4]byte{1}
		closest, err := n.ClosestPrecedingFinger(key)

		if err != nil {
			t.Error(err)
		}

		if !keyspace.Equal(closest.GetID(), n.GetID()) {
			t.Errorf("expected %x got %x", n.GetID(), closest.GetID())
		}
	})

	t.Run("other vals in ft", func(t *testing.T) {
		n.fingerTable[0] = newRemoteNodeWithID(newAddr(t), [4]byte{2})
		n.fingerTable[1] = newRemoteNodeWithID(newAddr(t), [4]byte{4})
		n.fingerTable[2] = newRemoteNodeWithID(newAddr(t), [4]byte{8})

		key := [4]byte{5}

		closest, err := n.ClosestPrecedingFinger(key)

		if err != nil {
			t.Error(err)
		}

		if !keyspace.Equal(closest.GetID(), [4]byte{4}) {
			t.Errorf("expected %x got %x", [4]byte{4}, closest.GetID())
		}
	})
}

func TestFindPredecessor(t *testing.T) {
	clearRemotes()
	t.Run("simple case - node only aware of itself", func(t *testing.T) {
		n := &Node{addr: Addr1, id: [4]byte{1}}
		n.predecessor = n.toRemoteNode()
		n.successor = n.toRemoteNode()
		key := [4]byte{0}
		nRemote := n.toRemoteNode()
		n.fingerTable = newFingerTable(nRemote)
		pred, err := n.FindPredecessor(key)

		if err != nil {
			t.Error(err)
		}

		if !keyspace.Equal(pred.GetID(), n.id) {
			t.Errorf("expected predecessor to be %x, got %x", n.id, pred.GetID())
		}
	})
	clearRemotes()
	t.Run("2 nodes - hi key", func(t *testing.T) {
		n1 := &Node{addr: newAddr(t), id: [4]byte{1}}

		n2 := &RemoteNodeMock{
			Addr:                   newAddr(t),
			ID:                     [4]byte{2},
			getSuccessorResponse:   n1.toRemoteNode(),
			getPredecessorResponse: n1.toRemoteNode(),
		}

		n1.predecessor = n2
		n1.successor = n2
		n1.fingerTable = newFingerTable(n2)

		key := [4]byte{3}

		p, err := n1.FindPredecessor(key)

		if err != nil {
			t.Error(err)
		}

		if !keyspace.Equal(p.GetID(), n2.GetID()) {
			t.Errorf("expected predecessor to be %x, got %x, %x", n2.GetID(), p.GetID(), n1.GetID())
		}
	})
	clearRemotes()
	t.Run("2 nodes - lo key", func(t *testing.T) {
		n1 := &Node{addr: newAddr(t), id: [4]byte{4}}

		n1Mock := &RemoteNodeMock{
			Addr: n1.addr,
			ID:   n1.id,
		}

		n2 := &RemoteNodeMock{
			Addr:                   newAddr(t),
			ID:                     [4]byte{1},
			getSuccessorResponse:   n1Mock,
			getPredecessorResponse: n1Mock,
		}

		n1.predecessor = n2
		n1.successor = n2
		n1.fingerTable = newFingerTable(n2)

		key := [4]byte{0}

		p, err := n1.FindPredecessor(key)

		if err != nil {
			t.Error(err)
		}

		if !keyspace.Equal(p.GetID(), n1.GetID()) {
			t.Errorf("expected predecessor to be % x, got % x", n1.GetID(), p.GetID())
		}
	})
}
