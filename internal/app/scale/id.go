package scale

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
)

// M keyspace
const M = 160

// Key 20 byte key
type Key = [M / 8]byte

// GenerateKey hash a string
func GenerateKey(str string) Key {
	h := sha1.New()
	h.Write([]byte(str))

	return ByteArrayToKey(h.Sum(nil))
}

// StringToKey convert a string directly to a 20 byte key
func StringToKey(str string) Key {
	return ByteArrayToKey([]byte(str))
}

// ByteArrayToKey convert a variable length byte array to
// 20 byte key
func ByteArrayToKey(arr []byte) Key {
	var key Key
	copy(key[:], arr)
	return key
}

// KeyToString Convert a key to a string
func KeyToString(key Key) string {
	return fmt.Sprintf("%x", key)
}

// Between returns whether n is between lower and upper
func Between(x, a, b Key) bool {

	X := x[:]
	A := a[:]
	B := b[:]
	// if A > B:
	if bytes.Compare(A, B) > 0 {
		// X is between A and B if X > A or X < B
		return bytes.Compare(X, A) > 0 || bytes.Compare(X, B) < 0
	}

	return bytes.Compare(A, X) < 0 && bytes.Compare(X, B) < 0
}

// BetweenRightInclusive returns whether n is between lower and upper, upper
// inclusive
func BetweenRightInclusive(x, a, b Key) bool {
	if bytes.Compare(a[:], b[:]) > 0 {
		return Between(x, a, b) || bytes.Equal(x[:], a[:])
	}
	return Between(x, a, b) || bytes.Equal(x[:], b[:])
}
