package scale

import "math/big"

/*
Some utility functions to help out with id hashes
*/

func hash

func IDToString(id []byte) string {
	key := big.Int{}
	key.SetBytes(id)

	return key.String()
}
