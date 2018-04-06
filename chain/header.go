package chain

import (
	"../utils"
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"
	"time"
)

const (
	MAX_UINT = ^uint(0)
	THREADS  = 8
)

type BlockHeader struct {
	Height              int
	PrevHash, BlockHash string
	Timestamp           string
	Data                []byte
	Nonce               uint
}

func NewBlockHeader(prev *BlockHeader, data []byte) (err error, block BlockHeader) {
	err = nil
	block = BlockHeader{Height: 1}
	if prev != nil {
		block.Height = prev.Height + 1
		block.PrevHash = prev.BlockHash
	}
	block.Timestamp = time.Now().Format(time.ANSIC)
	block.Data = data
	block.BlockHash = block.ComputeHash()
	return
}

func (hdr *BlockHeader) stringToHash() string {
	dataHash := utils.ComputeHashEncoded(hdr.Data)
	return hdr.Timestamp +
		hdr.PrevHash +
		dataHash +
		strconv.Itoa(int(hdr.Nonce))
}

func (hdr *BlockHeader) ComputeHash() string {
	return utils.ComputeHashEncoded([]byte(hdr.stringToHash()))
}

func (hdr *BlockHeader) Verify(prevHdr *BlockHeader) bool {
	if hdr.Height != prevHdr.Height+1 {
		return false
	}
	if hdr.PrevHash != prevHdr.BlockHash {
		return false
	}
	if prevHash := prevHdr.ComputeHash(); hdr.PrevHash != prevHash {
		return false
	}
	return true
}

func CheckHashOk(data []byte, difficulty uint) bool {
	for currentByte := 0; difficulty > 0; currentByte++ {
		if difficulty <= 8 {
			if bt := data[currentByte] >> (8 - difficulty); bt != 0 {
				return false
			}
			difficulty = 0
		} else {
			if data[currentByte] != 0 {
				return false
			}
			difficulty -= 8
		}
	}
	return true
}

func Mine(hdr *BlockHeader, difficulty uint, c <-chan uint, res_ch chan<- uint, wg *sync.WaitGroup) {
	defer func() {
		log.Println("Done")
	}()
	for {
		val, ok := <-c
		if !ok {
			break
		}
		hdr.Nonce = val
		hashStr := hdr.stringToHash()
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

func (hdr *BlockHeader) MineNext(difficulty uint, threads int) bool {
	now := time.Now().UnixNano()
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
		case _ = <-res_ch:
			break f_loop
		default:
			c <- i
		}
	}
	close(c)
	log.Println("Wait")
	wg.Wait()
	time_elapsed := float64(time.Now().UnixNano()-now) / math.Pow10(6)
	log.Printf("Time elapsed: %f ms", time_elapsed)
	log.Printf("Hashrate: %f", 1000*float64(hdr.Nonce)/time_elapsed)
	return true
}

func (hdr *BlockHeader) Print() string {
	dataHash := utils.ComputeHashEncoded(hdr.Data)
	return fmt.Sprintf("Height: %d, Timestamp: %s,\nBlockhash: %s, PrevBlockHash: %s\nData hash: %s\nNonce: %d",
		hdr.Height, hdr.Timestamp, hdr.BlockHash, hdr.PrevHash, dataHash, hdr.Nonce)
}
