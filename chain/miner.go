package chain

import (
	"../utils"
	"github.com/davecgh/go-spew/spew"
	"log"
	"sync"
	//"math"
)

const (
	THREADS = 3
)

type Miner interface {
	MineNext(hdr *BlockHeader, difficulty uint) bool
	GetResult() uint
}

type SimpleMiner struct {
	result uint
}

type MultiThreadMiner struct {
	SimpleMiner
}

type MultiThreadRangeMiner struct {
	SimpleMiner
}

func (miner *SimpleMiner) MineNext(hdr *BlockHeader, difficulty uint) bool {
	for i := uint(0); i < MAX_UINT; i++ {
		hdr.Nonce = i
		hashStr := hdr.stringToHash()
		hash := utils.ComputeHash([]byte(hashStr))
		if CheckHashOk(hash, difficulty) {
			hdr.BlockHash = utils.ComputeHashEncoded([]byte(hashStr))
			miner.result = i
			log.Printf("Found %d", miner.result)
			return true
			break
		}
	}
	return false
}

func (miner *SimpleMiner) GetResult() uint {
	log.Println("get result ", miner.result)
	return miner.result
}

func (miner *MultiThreadMiner) MineNext(hdr *BlockHeader, difficulty uint) bool {
	c := make(chan uint, THREADS*1024)
	res_ch := make(chan uint, THREADS)
	var wg sync.WaitGroup
	for i := 1; i <= THREADS; i++ {
		wg.Add(1)
		go Mine(hdr, difficulty, c, res_ch, &wg)
	}
f_loop:
	for i := uint(0); i < MAX_UINT; i++ {
		select {
		case nonce := <-res_ch:
			log.Println("nonce found: ", nonce)
			miner.result = nonce
			break f_loop
		default:
			c <- i
		}
	}
	close(c)
	wg.Wait()
	return true
}

func Mine(hdr *BlockHeader, difficulty uint, c <-chan uint, res_ch chan<- uint, wg *sync.WaitGroup) {
	defer func() {
		_ = <-c
	}()
	hdr_cpy := *hdr
	for {
		val, ok := <-c
		if !ok {
			break
		}
		hdr_cpy.Nonce = val
		hashStr := hdr_cpy.stringToHash()
		hash := utils.ComputeHash([]byte(hashStr))
		if CheckHashOk(hash, difficulty) {
			hdr.BlockHash = utils.ComputeHashEncoded([]byte(hashStr))
			log.Printf("Found %d", val)
			res_ch <- val
			break
		}
	}
	wg.Done()
}

type Range struct {
	from    uint
	to      uint
	current uint
}

func (rng *Range) Next() (val uint, ok bool) {
	if rng.current >= rng.to {
		return 0, false
	}
	rng.current++
	return rng.current - 1, true
}

func makeRange(from, to uint) Range {
	return Range{from, to, from}
}

func (miner *MultiThreadRangeMiner) MineNext(hdr *BlockHeader, difficulty uint) bool {
	var ranges [THREADS]Range
	rangePower := MAX_UINT / THREADS
	log.Println("Max int: ", MAX_UINT)
	for i := range ranges {
		ranges[i] = makeRange(uint(i)*rangePower, uint(i+1)*rangePower)
		log.Println(spew.Sdump(ranges[i]))
	}
	return true
}
