package scale

import (
	"fmt"

	"crypto/sha1"
)

/*
Some utility functions to help out with id hashes
*/

func generateHash(key string) []byte {
	h := sha1.New()
	h.Write([]byte(key))
	checksum := h.Sum(nil)
	fmt.Printf("%v\n", len(checksum))
	// hardcoded len at the moment, will eventually be passed by config.KeySize
	return checksum[:32]

}

// for debugging purposes
func IDToString(id []byte) string {
	return fmt.Sprintf("%x\n", id)
}

// we might need to pad ids since the lengths could be different
