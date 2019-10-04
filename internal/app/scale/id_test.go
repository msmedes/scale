package scale

import (
	"testing"
)

func TestBetween(t *testing.T) {
	t.Parallel()

	tests := []struct {
		n     string
		lower string
		upper string
		exp   bool
	}{
		{n: "10", lower: "0", upper: "20", exp: true},
		{n: "10", lower: "20", upper: "20", exp: false},
		{n: "529", lower: "527", upper: "789", exp: true},
		{n: "a", lower: "a", upper: "z", exp: false},
	}
	for _, test := range tests {
		want, got := test.exp, Between(
			StringToKey(test.n),
			StringToKey(test.lower),
			StringToKey(test.upper),
		)

		if got != want {
			t.Fatalf(
				"expected %t for between(%s, %s, %s), got %t",
				want,
				test.n,
				test.lower,
				test.upper,
				got,
			)
		}
	}
}
