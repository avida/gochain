package chain

import (
	"../utils"
  "sync"
	//"math"
)

const (
	THREADS = 8
  LOG_PREFIX = "miner"
)

type Miner interface {
	MineNext(hdr *BlockHeader, difficulty uint) bool
	GetResult() uint
  GetHashProcessed() uint
}

type SimpleMiner struct {
	result uint
  hashProcessed uint
}

type MultiThreadMiner struct {
	SimpleMiner
}

type MultiThreadRangeMiner struct {
	SimpleMiner
}

func (miner *SimpleMiner) MineNext(hdr *BlockHeader, difficulty uint) bool {
  logger := utils.GetLogger(LOG_PREFIX)
	for i := uint(0); i < MAX_UINT; i++ {
		hdr.Nonce = i
		hashStr := hdr.stringToHash()
		hash := utils.ComputeHash([]byte(hashStr))
		if CheckHashOk(hash, difficulty) {
			hdr.BlockHash = utils.ComputeHashEncoded([]byte(hashStr))
			miner.hashProcessed = i
      miner.result = i
			logger.Printf("Found %d", miner.result)
			return true
			break
		}
	}
	return false
}

func (miner *SimpleMiner) GetResult() uint {
	return miner.result
}

func (miner *SimpleMiner) GetHashProcessed() uint {
	return miner.hashProcessed
}

func (miner *MultiThreadMiner) MineNext(hdr *BlockHeader, difficulty uint) bool {
  logger := utils.GetLogger("miner")
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
			logger.Println("nonce found: ", nonce)
			miner.hashProcessed = nonce
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
  logger := utils.GetLogger(LOG_PREFIX)
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
			logger.Printf("Found %d", val)
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

func (rng *Range) Diff() uint {
  return rng.current - rng.from
}

func makeRange(from, to uint) Range {
	return Range{from, to, from}
}

func MineRange(hdr *BlockHeader, difficulty uint, rng *Range,
               res_ch chan<- uint, done_ch <-chan bool,
               wg *sync.WaitGroup) {
	defer func() {
    wg.Done()
	}()
  logger := utils.GetLogger(LOG_PREFIX)
	hdr_cpy := *hdr
  var i uint
  ok := true

  for {
  select {
    case _,_ = <-done_ch:
      return
    default:
      cntr:=0
      for ; ok && cntr < 10; i, ok = rng.Next(){
        hdr_cpy.Nonce = i
        hashstr := hdr_cpy.stringToHash()
        hash := utils.ComputeHash([]byte(hashstr))
        if CheckHashOk(hash, difficulty) {
          hdr.BlockHash = utils.ComputeHashEncoded([]byte(hashstr))
          logger.Printf("found %d", i)
          res_ch <- i
          return
        }
        cntr++
      }
  }
  }
}

func (miner *MultiThreadRangeMiner) MineNext(hdr *BlockHeader, difficulty uint) bool {
	var ranges [THREADS]Range
	var wg sync.WaitGroup
  logger := utils.GetLogger(LOG_PREFIX)
  res_ch := make(chan uint, THREADS)
  done_ch:= make(chan bool)
	rangePower := MAX_UINT / THREADS
  for i:= 0; i< THREADS; i++ {
    ranges[i] = makeRange(uint(i)*rangePower, uint(i+1)*rangePower)
  }
	for i := range ranges {
    wg.Add(1)
    go MineRange(hdr, difficulty, &ranges[i], res_ch, done_ch, &wg )
	}
  select {
    case nonce := <-res_ch:
      logger.Println("nonce found: ", nonce)
      miner.result = nonce
  }
  close(done_ch)
  wg.Wait()
  miner.hashProcessed = 0
  for i:=range ranges {
    miner.hashProcessed += ranges[i].Diff()
  }
	return true
}
