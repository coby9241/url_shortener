package utils

import (
	"crypto/sha256"
	"math/big"
)

const base62Alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// EncodeToBase62 encodes a string (treated as raw bytes) into a base62 string.
func EncodeToBase62(s string) string {
	num := new(big.Int).SetBytes([]byte(s))
	if num.Cmp(big.NewInt(0)) == 0 {
		return string(base62Alphabet[0])
	}

	var result []byte
	base := big.NewInt(62)
	n := new(big.Int).Set(num)
	mod := new(big.Int)

	for n.Cmp(big.NewInt(0)) > 0 {
		n.DivMod(n, base, mod)
		result = append(result, base62Alphabet[mod.Int64()])
	}

	// Reverse the result slice to get the correct order
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

// GenerateURLHash hashes a URL using SHA-256 and encodes the resulting hash to Base62.
// It returns a hash string truncated to the specified length.
func GenerateURLHash(url string, length int) string {
	hash := sha256.Sum256([]byte(url))
	base62Str := EncodeToBase62(string(hash[:]))

	if len(base62Str) > length {
		return base62Str[:length]
	}
	return base62Str
}
