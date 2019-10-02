package scale

import (
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
