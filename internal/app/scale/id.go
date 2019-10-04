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
func Between(n, lower, upper Key) bool {
	x := n[:]
	lo := lower[:]
	hi := upper[:]

	if bytes.Compare(lo, hi) > 0 {
		log.Fatalf("unexpected bounds: %s/%s", lo, hi)
	}

	return bytes.Compare(lo, x) < 0 && bytes.Compare(x, hi) < 0
}

// BetweenRightInclusive returns whether n is between lower and upper, upper
// inclusive
func BetweenRightInclusive(n, lower, upper Key) bool {
	return Between(n, lower, upper) || bytes.Equal(n[:], upper[:])
}
