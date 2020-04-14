package util

import (
	"crypto/rand"
	"math/big"
)

const (
	characters  = "bcdfghjklmnpqrstvwxz2456789"
	tokenLength = 54
)

var charsLength = big.NewInt(int64(len(characters)))

func GenerateRandomKey() (string, error) {
	key := make([]byte, tokenLength)
	for i := range key {
		r, err := rand.Int(rand.Reader, charsLength)
		if err != nil {
			return "", err
		}
		key[i] = characters[r.Int64()]
	}
	return string(key), nil
}
