package chain

import (
	"../utils"
	"sync"
	"time"
)

const (
	LOG_PREFIX = "miner"
)

type MinerResult struct {
	Nonce     uint
	Timestamp string
}

type Miner interface {
	MineNext(hdr *BlockHeader, difficulty uint) bool
	GetResult() MinerResult
	GetHashProcessed() uint
}

type SimpleMiner struct {
	result        MinerResult
	hashProcessed uint
}

type MultiThreadMiner struct {
	SimpleMiner
	Threads int
}

type RangeMiner struct {
	SimpleMiner
	Threads int
}

func (miner *SimpleMiner) MineNext(hdr *BlockHeader, difficulty uint) bool {
	logger := utils.GetLogger(LOG_PREFIX)
	for i := uint(0); i < MAX_UINT; i++ {
		currentTime := time.Now().Format(time.ANSIC)
		hdr.Timestamp = currentTime
		hdr.Nonce = i
		hashStr := hdr.stringToHash()
		hash := utils.ComputeHash([]byte(hashStr))
		if CheckHashOk(hash, difficulty) {
			hdr.BlockHash = utils.ComputeHashEncoded([]byte(hashStr))
			miner.hashProcessed = i
			miner.result = MinerResult{i, currentTime}
			logger.Printf("Found %d", miner.result)
			return true
			break
		}
	}
	return false
}

func (miner *SimpleMiner) GetResult() MinerResult {
	return miner.result
}

func (miner *SimpleMiner) GetHashProcessed() uint {
	return miner.hashProcessed
}

func (miner *MultiThreadMiner) MineNext(hdr *BlockHeader, difficulty uint) bool {
	logger := utils.GetLogger("miner")
	c := make(chan uint, miner.Threads*1024)
	res_ch := make(chan MinerResult, miner.Threads)
	var wg sync.WaitGroup
	for i := 1; i <= miner.Threads; i++ {
		wg.Add(1)
		go Mine(hdr, difficulty, c, res_ch, &wg)
	}
f_loop:
	for i := uint(0); i < MAX_UINT; i++ {
		select {
		case result := <-res_ch:
			logger.Println("nonce found: ", result.Nonce)
			miner.result = result
			miner.hashProcessed = result.Nonce
			break f_loop
		default:
			c <- i
		}
	}
	close(c)
	wg.Wait()
	return true
}

func Mine(hdr *BlockHeader, difficulty uint, c <-chan uint, res_ch chan<- MinerResult, wg *sync.WaitGroup) {
	defer func() {
		_ = <-c
	}()
	hdr_cpy := *hdr
	logger := utils.GetLogger(LOG_PREFIX)
	timeout := time.After(time.Second)
	currentTime := time.Now().Format(time.ANSIC)
	hdr_cpy.Timestamp = currentTime
for_loop:
	for {
		select {
		case val, ok := <-c:
			if !ok {
				break for_loop
			}
			hdr_cpy.Timestamp = currentTime
			hdr_cpy.Nonce = val
			hashStr := hdr_cpy.stringToHash()
			hash := utils.ComputeHash([]byte(hashStr))
			if CheckHashOk(hash, difficulty) {
				hdr.BlockHash = utils.ComputeHashEncoded([]byte(hashStr))
				logger.Printf("Found %d", val)
				res_ch <- MinerResult{val, currentTime}
				break for_loop
			}
		case <-timeout:
			timeout = time.After(time.Second)
			currentTime = time.Now().Format(time.ANSIC)
			hdr_cpy.Timestamp = currentTime
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

func MakeRangeMiner(threads int) *RangeMiner {
	var miner RangeMiner
	miner.Threads = threads
	return &miner
}
func MakeMultiThreadMiner(threads int) *MultiThreadMiner {
	var miner MultiThreadMiner
	miner.Threads = threads
	return &miner
}

func MineRange(hdr *BlockHeader, difficulty uint, rng *Range,
	res_ch chan<- MinerResult, done_ch <-chan bool,
	wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	logger := utils.GetLogger(LOG_PREFIX)
	hdr_cpy := *hdr
	var i uint
	ok := true
	currentTime := time.Now().Format(time.ANSIC)
	timeout := time.After(time.Second)
	res := MinerResult{}
	res.Timestamp = currentTime
	hdr_cpy.Timestamp = currentTime

	for {
		select {
		case _, _ = <-done_ch:
			return
		case <-timeout:
			// Update timestamp once per second
			currentTime = time.Now().Format(time.ANSIC)
			res.Timestamp = currentTime
			hdr_cpy.Timestamp = currentTime
			timeout = time.After(time.Second)
		default:
			cntr := 0
			for ; ok && cntr < 100; i, ok = rng.Next() {
				hdr_cpy.Nonce = i
				hashstr := hdr_cpy.stringToHash()
				hash := utils.ComputeHash([]byte(hashstr))
				if CheckHashOk(hash, difficulty) {
					hdr.BlockHash = utils.ComputeHashEncoded([]byte(hashstr))
					logger.Printf("found %d", i)
					res.Nonce = i
					res_ch <- res
					return
				}
				cntr++
			}
		}
	}
}

func (miner *RangeMiner) MineNext(hdr *BlockHeader, difficulty uint) bool {
	ranges := make([]Range, miner.Threads)
	var wg sync.WaitGroup
	logger := utils.GetLogger(LOG_PREFIX)
	res_ch := make(chan MinerResult, miner.Threads)
	done_ch := make(chan bool)
	rangePower := MAX_UINT / uint(miner.Threads)
	for i := 0; i < miner.Threads; i++ {
		ranges[i] = makeRange(uint(i)*rangePower, uint(i+1)*rangePower)
	}
	for i := range ranges {
		wg.Add(1)
		go MineRange(hdr, difficulty, &ranges[i], res_ch, done_ch, &wg)
	}
	select {
	case result := <-res_ch:
		logger.Printf("nonce found: %ui, ts: %s ", result.Nonce, result.Timestamp)
		miner.result = result
	}
	close(done_ch)
	wg.Wait()
	miner.hashProcessed = 0
	for i := range ranges {
		miner.hashProcessed += ranges[i].Diff()
	}
	return true
}
