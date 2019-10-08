package scale

import (
	"testing"
)

func TestBetween(t *testing.T) {
	t.Parallel()

	tests := []struct {
		x   string
		a   string
		b   string
		exp bool
	}{
		{x: "10", a: "0", b: "20", exp: true},
		{x: "10", a: "20", b: "20", exp: false},
		{x: "529", a: "527", b: "789", exp: true},
		{x: "a", a: "a", b: "z", exp: false},

		{x: "20", a: "527", b: "277", exp: true},
		{x: "788", a: "527", b: "277", exp: true},
		{x: "20", a: "5", b: "2", exp: false},
		{x: "1", a: "5", b: "2", exp: true},
		{x: "3", a: "5", b: "2", exp: false},
		{x: "20", a: "2", b: "5", exp: true},
	}
	for _, test := range tests {
		want, got := test.exp, Between(
			StringToKey(test.x),
			StringToKey(test.a),
			StringToKey(test.b),
		)

		if got != want {
			t.Fatalf(
				"expected %t for between(%s, %s, %s), got %t",
				want,
				test.x,
				test.a,
				test.b,
				got,
			)
		}
	}
}
