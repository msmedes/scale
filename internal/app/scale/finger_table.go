package scale

import (
	"math/big"
)

type FingerTable []*finger

type finger struct {
	Id         Key
	RemoteNode *Node
}

func NewFingerTable(m int) FingerTable {
	return make(FingerTable, m)
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
