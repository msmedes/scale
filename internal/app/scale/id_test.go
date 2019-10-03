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

		// when a < b the number should be between

		{x: "10", a: "0", b: "20", exp: true},
		{x: "10", a: "20", b: "20", exp: true},
		{x: "10", a: "12", b: "10", exp: false},
		{x: "529", a: "527", b: "789", exp: true},
		{x: "a", a: "a", b: "z", exp: true},
	}
	for _, test := range tests {
		x := GenerateKey(test.x)
		a := GenerateKey(test.a)
		b := GenerateKey(test.b)
		want, got := test.exp, between(x, a, b)
		if got != want {
			t.Fatalf("expected %t for between(%s, %s, %s), got %t", want, test.x, test.a, test.b, got)
		}
	}
}
