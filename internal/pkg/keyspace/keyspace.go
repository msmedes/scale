package keyspace

import (
	"bytes"
	"crypto/sha1"
	"fmt"
  "github.com/msmedes/scale/internal/pkg/scale"
)


// Equal whether 2 keys are equal
func Equal(a scale.Key, b scale.Key) bool {
	return bytes.Equal(a[:], b[:])
}

// GenerateKey hash a string

func GenerateKey(str string) Key {
	hash := sha1.Sum([]byte(str))
	return ByteArrayToKey(hash[:])
}

// StringToKey convert a string directly to a 20 byte key
func StringToKey(str string) scale.Key {
	return ByteArrayToKey([]byte(str))
}

// ByteArrayToKey convert a variable length byte array to
// 20 byte key
func ByteArrayToKey(arr []byte) scale.Key {
	var key scale.Key
	copy(key[:], arr)
	return key
}

// KeyToString Convert a key to a string
func KeyToString(key scale.Key) string {
	return fmt.Sprintf("%x", key)
}

// Between returns whether x is between a and b
func Between(x, a, b scale.Key) bool {
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
func BetweenRightInclusive(x, a, b scale.Key) bool {
	if bytes.Compare(a[:], b[:]) > 0 {
		return Between(x, a, b) || bytes.Equal(x[:], a[:])
	}
	return Between(x, a, b) || bytes.Equal(x[:], b[:])
}
