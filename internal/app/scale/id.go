package scale

import (
	"bytes"
	"crypto/sha1"
	"fmt"
)

const M = 160

type Key = [M / 8]byte

func GenerateKey(str string) Key {
	h := sha1.New()
	h.Write([]byte(str))

	return ByteArrayToKey(h.Sum(nil))
}

func StringToKey(str string) Key {
	return ByteArrayToKey([]byte(str))
}

func ByteArrayToKey(arr []byte) Key {
	var key Key
	copy(key[:], arr)
	return key
}

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
