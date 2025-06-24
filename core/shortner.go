package core

import (
	"crypto/rand"
	"math/big"
)

func GenerateShortCode(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	r := make([]byte, length)
	for i := range r {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err

		}
		r[i] = charset[n.Int64()]
	}
	return string(r), nil
}
