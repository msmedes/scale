package scale

import (
	"bytes"
	"fmt"
	"math/big"
)

type FingerTable []*finger

type finger struct {
	ID         Key
	RemoteNode *Node
}

func NewFingerTable(m int, n *Node) FingerTable {
	ft := make([]*finger, m)
	for i := range ft {
		ft[i] = newFinger(n.ID, n)
	}
	return ft
}

// func InitFingerTable(n *Node) {
// 	// when a new chord is created the first node should just have a finger
// 	// table filled with fingers pointing to itself since it doesn't have
// 	// a predecessor or successor yet
// 	for i := range n.fingerTable {
// 		n.fingerTable[i] = newFinger(n.ID, n)
// 	}
// }

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

func (ft FingerTable) String() string {
	var buf bytes.Buffer

	for _, val := range ft {
		buf.WriteString(fmt.Sprintf(
			"\n{id:%v\tnodeId:%v",
			IdToString(val.ID),
			IdToString(val.RemoteNode.ID),
		))
	}
	return buf.String()
}
