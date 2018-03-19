package main

import (
	"./utils"
	"crypto/dsa"
	"crypto/rand"
	"log"
	"os"
)

const FirstData = "This is first block data"

func main() {
	defer func() {
		if r := recover(); r != nil {
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
		hash := utils.ComputeHash([]byte(FirstData))
		log.Println(hash)
	default:
		log.Println("wrong argument: ", os.Args[1])
	}
}
