package main

import (
	"./chain"
	"./db"
	"./utils"
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"log"
	"net/http"
	"time"
)

var main_logger *log.Logger

type BlockExplorer struct {
	Ledger *chain.Chain
}

func (explorer *BlockExplorer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<ul>")
	for _, block := range *explorer.Ledger {
		if block.Nonce == 0 {
			continue
		}
		fmt.Fprintf(w, "<li>%s</li>", block.Print())
	}
	fmt.Fprintf(w, "</ul>")
}

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
	main_logger.Println(db.ConnStr())

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
	explorerHandler := BlockExplorer{
		Ledger: &ledger,
	}

	s := http.Server{
		Addr:           fmt.Sprintf(":%d", *port),
		Handler:        &explorerHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1024,
	}

	go func() {
		main_logger.Fatal(s.ListenAndServe())
	}()

	miner := chain.MakeRangeMiner(*threads)
	main_logger.Printf("Mining with %d threads started, difficulty is %d", *threads, *difficulty)
	for {
		block := ledger.AddBlock()
		if block.MineNext(*difficulty, miner) {
			main_logger.Println("Block mined: %v", spew.Sdump(block))
		}
	}
}
