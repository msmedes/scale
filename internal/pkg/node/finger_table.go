package node

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/msmedes/scale/internal/pkg/keyspace"
	"github.com/msmedes/scale/internal/pkg/scale"
)

// Table contains nodes in network for lookups
type Table []*RemoteNode

// Finger finger
// type Finger struct {
// 	ID   scale.Key
// 	Addr string
// }

// NewFingerTable create and populate a finger table
func NewFingerTable(n *Node) Table {
	ft := make([]*scale.RemoteNode, scale.M)

	for i := range ft {
		ft[i] = NewRemoteNode(n.Addr)
	}

	return ft
}

// Math fingermath
func Math(n []byte, i int, m int) []byte {
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

func (f Finger) String() string {
	return fmt.Sprintf("%s", keyspace.KeyToString(f.ID))
}

func (ft Table) String() string {
	var buf bytes.Buffer

	buf.WriteString("\n")

	for _, val := range ft {
		str := fmt.Sprintf("%s\n", val.String())
		buf.WriteString(str)
	}

	return buf.String()
}
