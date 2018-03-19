package utils

import (
	"crypto/sha256"
	"encoding/base64"
	pseudorand "math/rand"
)

func ComputeHash(data []byte) string {
	result := sha256.Sum256(data)
	return base64.StdEncoding.EncodeToString(result[:])
}

func ReadRandom(len int) (result []byte, err error) {
	result = make([]byte, len, len)
	_, err = pseudorand.Read(result)
	return
}
