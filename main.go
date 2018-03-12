package main

import (
	"crypto/dsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"os"
)

const FirstData = "This is first block data"

type BlockHeader struct {
	Height             int
	PrevHash, DataHash string
	Timestamp          string
}

func computeHash(data string) string {
	result := sha256.Sum256([]byte(data))
	return base64.StdEncoding.EncodeToString(result[:])
}

func main() {
  defer func() {
    if r:= recover(); r != nil {
      log.Println("Unexpected error: ", r)
    }
  }()

	log.Println("This is my first golang application")
	switch os.Args[1] {
	case "key":
		params := new(dsa.Parameters)
		if err := dsa.GenerateParameters(params, rand.Reader, dsa.L1024N160); err != nil {
			log.Println("err")
			os.Exit(1)
		}
		private_key := new(dsa.PrivateKey)
		private_key.PublicKey.Parameters = *params
		dsa.GenerateKey(private_key, rand.Reader)
		log.Println(*private_key)
	case "hash":
		hash := computeHash(FirstData)
		log.Println(hash)
  default:
    log.Println("wrong argument: ", os.Args[1])
	}
}
