package node

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/msmedes/scale/internal/pkg/scale"
)

// Table contains nodes in network for lookups
type table []scale.RemoteNode

// Great news, there are now two ways we need to initialize
// a finger table, this is for when there is only one node
// in the chord
func newFingerTable(node scale.RemoteNode) table {
	ft := make([]scale.RemoteNode, scale.M)

	for i := range ft {
		ft[i] = node
	}

	return ft
}

// FingerMath fingermath
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

func (ft table) String() string {
	var buf bytes.Buffer

	buf.WriteString("\n")

	for _, val := range ft {
		str := fmt.Sprintf("%x", val.GetID())
		buf.WriteString(str)
		buf.WriteString("\n")
	}

	return buf.String()
}
