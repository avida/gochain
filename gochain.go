package main

import (
	"./chain"
	"./utils"
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"log"
	"net/http"
	"time"
)

const FirstData = "This is first block data"

var main_logger *log.Logger

func SetupLoggers() {
	utils.SetupLogger("miner")
	utils.SetupLogger("chain")
	utils.SetupLogger("header")
	utils.SetupLogger("test")
	main_logger = utils.SetupLogger("main")
	//utils.SetOutput("miner", utils.StdOut)
	//utils.SetOutput("header", utils.StdOut)
	//utils.SetOutput("chain", utils.StdOut)
	//utils.SetOutput("test", utils.StdOut)
	utils.SetOutput("main", utils.StdOut)
}

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
	SetupLoggers()
	threads := flag.Int("t", 8, "Number of threads")
	difficulty := flag.Uint("d", 24, "Hash Difficulty")
	port := flag.Uint("port", 8080, "Http port to listen")
	flag.Parse()

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
