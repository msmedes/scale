package scale

import (
	"bytes"
	"fmt"
	"math/big"
)

// FingerTable contains nodes in network for lookups
type FingerTable []*finger

type finger struct {
	ID         Key
	RemoteNode *Node
}

// NewFingerTable create and populate a finger table
func NewFingerTable(m int, n *Node) FingerTable {
	ft := make([]*finger, m)

	for i := range ft {
		ft[i] = newFinger(n.ID, n)
	}

	return ft
}

func newFinger(id Key, n *Node) *finger {
	return &finger{ID: id, RemoteNode: n}
}

func fingerMath(n []byte, i int, m int) []byte {
	twoExp := big.NewInt(2)
	twoExp.Exp(twoExp, big.NewInt(int64(i)), nil)
	mExp := big.NewInt(2)
	mExp.Exp(mExp, big.NewInt(int64(m)), nil)

	res := &big.Int{}
	res.SetBytes(n[:])
	res.Add(res, twoExp)
	res.Mod(res, mExp)

	return res.Bytes()
}

func (f finger) String() string {
	return fmt.Sprintf("%s", KeyToString(f.ID))
}

func (ft FingerTable) String() string {
	var buf bytes.Buffer

	buf.WriteString("\n")

	for _, val := range ft {
		str := fmt.Sprintf("%s\n", val.String())
		buf.WriteString(str)
	}

	return buf.String()
}
