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
	ID   []byte
	Node *Node
}

func (f *finger) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("id: %v node: %v", f.ID, f.Node.ID))

	return buf.String()
}

func newFingerTable(n *Node, m int) FingerTable {
	// TODO: A few finger table should take a node and an integer m
	// and construct a new finger table of fingers length m.
	fingerTable := make([]*finger, 8)
	for i := range fingerTable {
		fingerTable[i] = newFinger(fingerMath(n.ID, i, 8), n)
	}
	return fingerTable
}

// calculates the offset for the finger table through (n + 2^i) mod (2 ^ m)
func fingerMath(n []byte, i int, m int) []byte {
	twoExp := big.NewInt(2)
	twoExp.Exp(twoExp, big.NewInt(int64(i)), nil) // 8 is just a placeholder for the accuracy, eventually will be passed by config
	mExp := big.NewInt(2)
	mExp.Exp(mExp, big.NewInt(int64(m)), nil)

	res := &big.Int{}
	res.SetBytes(n)      // set the byte array to big int
	res.Add(res, twoExp) // do the addition of the first term
	res.Mod(res, mExp)   // mod that bad boi

	return res.Bytes()
}

func newFinger(id []byte, n *Node) *finger {
	return &finger{
		ID:   id,
		Node: n,
	}
}
