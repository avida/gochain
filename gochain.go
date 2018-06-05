package main

import (
	"./chain"
	"./utils"
	"./web"
	"flag"
	"log"
)

var main_logger *log.Logger

func main() {
	defer func() {
		if r := recover(); r != nil {
			main_logger.Println("Unexpected error: ", r)
		}
	}()

	utils.SetupLoggers()
	main_logger = utils.GetLogger("main")
	threads := flag.Int("t", 8, "Number of threads")
	difficulty := flag.Uint("d", 24, "Hash Difficulty")
	port := flag.Uint("port", 8080, "Http port to listen")
	configPath := flag.String("config", "test", "Path to config directory")
	flag.Parse()
	utils.ReadConf(*configPath)

	if len(flag.Args()) != 0 {
		switch flag.Args()[0] {
		case "config":
			main_logger.Printf("config ")
		case "rpc":
			main_logger.Printf("rpc")
		}
		main_logger.Printf("args: %v", flag.Args())
	}

	var ledger chain.Chain
	web.RunServer(&ledger, *port)

	miner := chain.MakeRangeMiner(*threads)
	main_logger.Printf("Mining with %d threads started, difficulty is %d", *threads, *difficulty)
	for {
		block := ledger.AddBlock()
		if block.MineNext(*difficulty, miner) {
			main_logger.Printf("Block w/ hash %s mined, height is %d", block.BlockHash, block.Height)
		}
	}
}
