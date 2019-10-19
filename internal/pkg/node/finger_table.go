package node

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/msmedes/scale/internal/pkg/scale"
)

// Table contains nodes in network for lookups
type table []scale.RemoteNode

func newFingerTable(node scale.RemoteNode) table {
	ft := make([]scale.RemoteNode, scale.M)

	for i := range ft {
		ft[i] = node
	}

	return ft
}

// FingerMath fingermath
func fingerMath(n []byte, i int) []byte {
	iInt := big.NewInt(2)
	iInt.Exp(iInt, big.NewInt(int64(i)), nil)
	mInt := big.NewInt(2)
	mInt.Exp(mInt, big.NewInt(int64(scale.M)), nil)

	res := &big.Int{} // res will pretty much be an accumulator
	res.SetBytes(n[:]).Add(res, iInt).Mod(res, mInt)

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
