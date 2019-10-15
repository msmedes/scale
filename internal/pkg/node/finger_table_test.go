package node

import (
	"bytes"
	"math/big"
	"testing"
)

func TestFingerMath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		n   int64
		i   int
		exp int64
	}{
		{n: 0, i: 0, exp: 1},
		{n: 2, i: 0, exp: 3},
		{n: 4, i: 0, exp: 5},
		{n: 8, i: 0, exp: 9},
		{n: 64, i: 0, exp: 65},
		{n: 256, i: 0, exp: 1},
		{n: 1000, i: 0, exp: 233},
		{n: 65563, i: 0, exp: 28},

		{n: 0, i: 2, exp: 4},
		{n: 2, i: 2, exp: 6},
		{n: 4, i: 2, exp: 8},
		{n: 8, i: 2, exp: 12},
		{n: 64, i: 2, exp: 68},
		{n: 256, i: 2, exp: 4},
		{n: 1000, i: 2, exp: 236},
		{n: 65563, i: 2, exp: 31},

		{n: 0, i: 8, exp: 0},
		{n: 2, i: 8, exp: 2},
		{n: 4, i: 8, exp: 4},
		{n: 8, i: 8, exp: 8},
		{n: 64, i: 8, exp: 64},
		{n: 256, i: 8, exp: 0},
		{n: 1000, i: 8, exp: 232},
		{n: 65563, i: 8, exp: 27},
	}

	for _, tt := range tests {
		n := big.NewInt(tt.n).Bytes()
		got := FingerMath(n, tt.i, 8)
		want := big.NewInt(tt.exp).Bytes()

		if !bytes.Equal(got[:], want) {
			t.Fatalf("expected %v, got %v for %v", tt.exp, (&big.Int{}).SetBytes(got[:]), tt)
		}
	}
}
