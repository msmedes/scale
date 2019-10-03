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

func ByteArrayToKey(arr []byte) Key {
	var key Key
	copy(key[:], arr)
	return key
}

func IdToString(id Key) string {
	return fmt.Sprintf("%x", id)
}

func between(x, a, b Key) bool {
	// return bytes.Compare(a[:], x[:]) == -1 && bytes.Compare(x[:], b[:]) == -1
	aSlice := a[:]
	bSlice := b[:]
	xSlice := x[:]

	switch bytes.Compare(aSlice, bSlice) {
	case -1:
		return bytes.Compare(aSlice, xSlice) == -1 && bytes.Compare(xSlice, bSlice) == -1
	case 1:
		return bytes.Compare(xSlice, aSlice) == -1 || bytes.Compare(xSlice, bSlice) == -1
	case 0:
		return bytes.Compare(xSlice, aSlice) != 0
	}

	return false
}

func betweenRightInclusive(x, a, b Key) bool {

	return between(x, a, b) || bytes.Equal(x[:], b[:])
}
