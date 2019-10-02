package scale

import (
	"bytes"
	"math/big"
	"testing"
)

func TestFingerMath(t *testing.T) {
	tests := []struct {
		n        int64
		i        int
		expected int64
	}{
		{n: 0, i: 0, expected: 1},
		{n: 2, i: 0, expected: 3},
		{n: 4, i: 0, expected: 5},
		{n: 8, i: 0, expected: 9},
		{n: 64, i: 0, expected: 65},
		{n: 256, i: 0, expected: 1},
		{n: 1000, i: 0, expected: 233},
		{n: 65563, i: 0, expected: 28},

		{n: 0, i: 2, expected: 4},
		{n: 2, i: 2, expected: 6},
		{n: 4, i: 2, expected: 8},
		{n: 8, i: 2, expected: 12},
		{n: 64, i: 2, expected: 68},
		{n: 256, i: 2, expected: 4},
		{n: 1000, i: 2, expected: 236},
		{n: 65563, i: 2, expected: 31},

		{n: 0, i: 8, expected: 0},
		{n: 2, i: 8, expected: 2},
		{n: 4, i: 8, expected: 4},
		{n: 8, i: 8, expected: 8},
		{n: 64, i: 8, expected: 64},
		{n: 256, i: 8, expected: 0},
		{n: 1000, i: 8, expected: 232},
		{n: 65563, i: 8, expected: 27},
	}

	for _, tt := range tests {
		n := big.NewInt(tt.n).Bytes()
		got := fingerMath(n, tt.i, 8)
		want := big.NewInt(tt.expected).Bytes()

		if !bytes.Equal(got[:], want) {
			t.Fatalf("expected %v, got %v for %v", tt.expected, (&big.Int{}).SetBytes(got[:]), tt)
		}
	}
}
