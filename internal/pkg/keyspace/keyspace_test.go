package keyspace

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
		{x: "20", a: "0", b: "10", exp: true},
		{x: "10", a: "20", b: "20", exp: false},
		{x: "529", a: "527", b: "789", exp: true},
		{x: "a", a: "a", b: "z", exp: false},

		{x: "20", a: "527", b: "277", exp: true},
		{x: "788", a: "277", b: "527", exp: true},
		{x: "3000", a: "3002", b: "3001", exp: true},
		{x: "5", a: "20", b: "2", exp: true},
		{x: "1", a: "2", b: "5", exp: true},
		{x: "3", a: "5", b: "2", exp: false},
	}
	for _, test := range tests {
		want, got := test.exp, Between(
			GenerateKey(test.x),
			GenerateKey(test.a),
			GenerateKey(test.b),
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

func TestBetweenRightInclusive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		x   string
		a   string
		b   string
		exp bool
	}{
		{x: "20", a: "0", b: "10", exp: true},
		{x: "10", a: "20", b: "20", exp: false},
		{x: "20", a: "20", b: "20", exp: true},
		{x: "20", a: "0", b: "20", exp: true},
		{x: "529", a: "527", b: "789", exp: true},
		{x: "a", a: "a", b: "z", exp: false},

		{x: "20", a: "527", b: "277", exp: true},
		{x: "788", a: "277", b: "527", exp: true},
		{x: "3000", a: "3002", b: "3001", exp: true},
		{x: "5", a: "20", b: "5", exp: true},
		{x: "1", a: "2", b: "5", exp: true},
		{x: "4", a: "5", b: "4", exp: true},
	}

	for _, test := range tests {
		want, got := test.exp, BetweenRightInclusive(
			GenerateKey(test.x),
			GenerateKey(test.a),
			GenerateKey(test.b),
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
