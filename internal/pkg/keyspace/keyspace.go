package keyspace

import (
	"bytes"
	"crypto/sha1"
	"fmt"

	"github.com/msmedes/scale/internal/pkg/scale"
)

// Equal whether 2 keys are equal
func Equal(a, b scale.Key) bool {
	return bytes.Equal(a[:], b[:])
}

// GenerateKey hash a string
func GenerateKey(str string) scale.Key {
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
	return fmt.Sprintf("% x", key)
}

// Between returns whether x is between a and b
func Between(x, a, b scale.Key) bool {
	// if A > B:
	if GT(a, b) {
		// X is between A and B if X > A or X < B
		return GT(x, a) || LT(x, b)
	}

	return GT(x, a) && LT(x, b)
}

// GT x greater than a
func GT(x, a scale.Key) bool {
	return bytes.Compare(x[:], a[:]) > 0
}

// LT x less than a
func LT(x, a scale.Key) bool {
	return bytes.Compare(x[:], a[:]) < 0
}

// GTE x greater than equal to a
func GTE(x, a scale.Key) bool {
	return bytes.Compare(x[:], a[:]) >= 0
}

// LTE x less than equal to a
func LTE(x, a scale.Key) bool {
	return bytes.Compare(x[:], a[:]) <= 0
}

// BetweenRightInclusive returns whether n is between lower and upper, upper
// inclusive
func BetweenRightInclusive(x, a, b scale.Key) bool {
	return Between(x, a, b) || Equal(x, b)
}
