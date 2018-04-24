package main

import (
	"./utils"
	"crypto/dsa"
	"crypto/rand"
	"flag"
	"log"
	"os"
)

const FirstData = "This is first block data"

var logger *log.Logger

func SetupLoggers() {
	utils.SetupLogger("miner")
	utils.SetupLogger("chain")
	utils.SetupLogger("header")
	utils.SetupLogger("test")
	logger = utils.SetupLogger("main")
	//utils.SetOutput("miner", utils.StdOut)
	//utils.SetOutput("chain", utils.StdOut)
	//utils.SetOutput("header", utils.File)
	//utils.SetOutput("test", utils.StdOut)
	utils.SetOutput("main", utils.StdOut)
}

func GenerateLedger() {
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			logger.Println("Unexpected error: ", r)
		}
	}()
	SetupLoggers()
	flag.Parse()

	logger.Println("This is my first golang application")
	switch os.Args[1] {
	case "key":
		params := new(dsa.Parameters)
		if err := dsa.GenerateParameters(params, rand.Reader, dsa.L1024N160); err != nil {
			logger.Println("err")
			os.Exit(1)
		}
		private_key := new(dsa.PrivateKey)
		private_key.PublicKey.Parameters = *params
		dsa.GenerateKey(private_key, rand.Reader)
		logger.Println(*private_key)
	case "hash":
		hash := utils.ComputeHash([]byte(FirstData))
		logger.Println(hash)
	default:
		logger.Println("wrong argument: ", os.Args[1])
	}
}
