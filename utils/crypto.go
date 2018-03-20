package utils

import (
	"crypto/sha256"
	"encoding/base64"
	pseudorand "math/rand"
)

func ComputeHash(data []byte)[]byte{
  res := sha256.Sum256(data)
  return res[:]
}

func ComputeHashEncoded(data []byte) string {
	hash := ComputeHash(data)
	return base64.StdEncoding.EncodeToString(hash)
}

func ReadRandom(len int) (result []byte, err error) {
	result = make([]byte, len, len)
	_, err = pseudorand.Read(result)
	return
}
