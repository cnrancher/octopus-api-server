package util

import (
	"crypto/rand"
	"math/big"
	rnd "math/rand"
	"strings"
	"time"
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

func GenerateRandomTempKey(length int) string {
	r := rnd.New(rnd.NewSource(time.Now().Unix()))
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return strings.ToLower(string(bytes))
}
