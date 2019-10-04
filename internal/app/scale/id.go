package scale

import (
	"bytes"
	"crypto/sha1"
	"fmt"
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

func between(n, lower, upper Key) bool {
	x := n[:]
	lo := lower[:]
	hi := upper[:]

	switch bytes.Compare(lo, hi) {
	case -1:
		return bytes.Compare(lo, x) <= 0 && bytes.Compare(x, hi) < 0
	case 1:
		return bytes.Compare(hi, x) < 0 && bytes.Compare(x, lo) <= 0
	}

	return false
}

func betweenRightInclusive(n, lower, upper Key) bool {
	return between(n, lower, upper) || bytes.Equal(n[:], upper[:])
}
