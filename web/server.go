package web

import (
	"../chain"
	"../utils"
	"fmt"
	"net/http"
	"time"
)

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

func RunServer(ledger *chain.Chain, port uint) {
	explorerHandler := BlockExplorer{
		Ledger: ledger,
	}

	s := http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        &explorerHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1024,
	}

	go func() {
		logger := utils.GetLogger("main")
		logger.Fatal(s.ListenAndServe())
	}()
}
