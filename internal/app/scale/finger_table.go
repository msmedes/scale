package scale

import (
	"bytes"
	"fmt"
	"math/big"
)

// FingerTable contains nodes in network for lookups
type FingerTable []*finger

type finger struct {
	ID Key
}

// NewFingerTable create and populate a finger table
func NewFingerTable(m int, ID Key) FingerTable {
	ft := make([]*finger, m)

	for i := range ft {
		ft[i] = &finger{ID: ID}
	}

	return ft
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

func (node *Node) fixNextFinger(next int) int {
	nextHash := fingerMath(node.ID[:], next, M)
	successor := node.findSuccessor(ByteArrayToKey(nextHash))
	finger := &finger{ID: successor.ID}
	node.fingerTable[next] = finger
	return next + 1
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
