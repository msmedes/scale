package node

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/msmedes/scale/internal/pkg/scale"
)

// Table contains nodes in network for lookups
type table []scale.RemoteNode

// FingerTable is a hell scape wherein up is down and down is up,
// the innocent die in ditches like dogs to the symphonic cackling
// of The FingerMath.  Here you too shall perish, driven to insanity by
// its machinations.  Abandon hope, all ye who entre here.
type FingerTable struct {
	Table  table
	starts map[int][]byte
}

func newFingerTable(node scale.RemoteNode) *FingerTable {
	ft := &FingerTable{
		Table:  make(table, scale.M),
		starts: make(map[int][]byte, scale.M),
	}

	for i := range ft.Table {
		ft.Table[i] = node
		id := node.GetID()
		ft.starts[i] = fingerMath(id[:], i)
	}

	// for i := range ft.Table {
	// 	keyStart := keyspace.KeyToString(keyspace.ByteArrayToKey(ft.starts[i]))
	// 	finger := keyspace.KeyToString(ft.Table[i].GetID())
	// 	fmt.Printf("%d start: %s, finger: %s\n", i, keyStart, finger)
	// }

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

	byteArray := res.Bytes()

	// this hack is bad and makes me feel bad
	keyLength := scale.M / 8
	if len(byteArray) < keyLength {
		for i := 0; i <= keyLength-len(byteArray); i++ {
			byteArray = append([]byte{0}, byteArray...)
		}
	}

	return byteArray
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
