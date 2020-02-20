package node

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/msmedes/scale/internal/pkg/keyspace"
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
		{n: 256, i: 0, exp: 257},
		{n: 1000, i: 0, exp: 1001},
		{n: 65563, i: 0, exp: 65564},
		{n: 4294967296, i: 0, exp: 1},

		{n: 0, i: 2, exp: 4},
		{n: 2, i: 2, exp: 6},
		{n: 4, i: 2, exp: 8},
		{n: 8, i: 2, exp: 12},
		{n: 64, i: 2, exp: 68},
		{n: 256, i: 2, exp: 260},
		{n: 1000, i: 2, exp: 1004},
		{n: 65563, i: 2, exp: 65567},
		{n: 4294967296, i: 2, exp: 4},

		{n: 0, i: 8, exp: 256},
		{n: 2, i: 8, exp: 258},
		{n: 4, i: 8, exp: 260},
		{n: 8, i: 8, exp: 264},
		{n: 64, i: 8, exp: 320},
		{n: 256, i: 8, exp: 512},
		{n: 1000, i: 8, exp: 1256},
		{n: 65563, i: 8, exp: 65819},

		{n: 0, i: 31, exp: 2147483648},
		{n: 2, i: 31, exp: 2147483650},
		{n: 4, i: 31, exp: 2147483652},
		{n: 8, i: 31, exp: 2147483656},
		{n: 64, i: 31, exp: 2147483712},
		{n: 256, i: 31, exp: 2147483904},
		{n: 1000, i: 31, exp: 2147484648},
		{n: 65563, i: 31, exp: 2147549211},
		{n: 4294967296, i: 31, exp: 2147483648},
		{n: 4294967295, i: 31, exp: 2147483647},
	}

	for _, tt := range tests {
		n := big.NewInt(tt.n).Bytes()
		got := keyspace.ByteArrayToKey(fingerMath(n, tt.i))
		expPrepend := createExp(got[:], tt.exp)
		exp := keyspace.ByteArrayToKey(expPrepend)

		if !bytes.Equal(got[:], exp[:]) {
			t.Fatalf("expected %v for n %v, got %v for %+v %v", exp, n, got, tt, expPrepend)
		}
	}
}

func createExp(got []byte, exp int64) []byte {
	expBig := big.NewInt(exp).Bytes()
	if len(expBig) < len(got) {
		for i := 0; i <= len(got)-len(expBig); i++ {
			expBig = append([]byte{0}, expBig...)
		}
	}
	return expBig
}
