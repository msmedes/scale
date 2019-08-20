package scale

import (
	"bytes"
	"fmt"
	"math/big"
)

/*
(n + 2^i mod (2^m) is the formula to remember for finger tables/chord
This will be used to calculate the offset for each entry of the finger table
from the root node.
n = number of nodes in the chord
i = the ith entry in the finger table
m = number of bits in the hash (I believe most implementations use 160 bit
sha1)
*/

type FingerTable []*finger

type finger struct {
	ID         []byte
	RemoteNode *Node
}

func (f *finger) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("id: %v node: %v", f.ID, f.RemoteNode.ID))

	return buf.String()
}

func newFingerTable(n *Node, m int) FingerTable {
	fingerTable := make([]*finger, 8)
	for i := range fingerTable {
		fingerTable[i] = newFinger(fingerMath(n.ID, i, 8), n)
	}
	return fingerTable
}

// calculates the offset for the finger table through (n + 2^i) mod (2 ^ m)
func fingerMath(n []byte, i int, m int) []byte {
	twoExp := big.NewInt(2)
	// nil is for "accuracy", for some reason before it was giving me zeros
	// instead of big ints, not great for doing math
	twoExp.Exp(twoExp, big.NewInt(int64(i)), nil)
	mExp := big.NewInt(2)
	mExp.Exp(mExp, big.NewInt(int64(m)), nil)

	res := &big.Int{}
	res.SetBytes(n)      // set the byte array to big int
	res.Add(res, twoExp) // do the addition of the first term
	res.Mod(res, mExp)   // mod that bad boi

	return res.Bytes()
}

func (n *Node) fixFingerTable() {
	/*
	  for i in range m (the key size):
	  use fingerMath to calculate where we should look next
	  use an RPC call to find the Successor of this node based on that hash
	  create a new finger pointing to the Successor node
	  Lock the fingertables mutex (idk can we use channels for this? it
	  could theoretically block? like each index listens on a channel? idk
	  that might be a lot when we could just lock and unlock)
	  update the finger table at i
	  unlock
	*/
}

func newFinger(id []byte, remoteNode *Node) *finger {
	return &finger{
		ID:         id,
		RemoteNode: remoteNode,
	}
}
