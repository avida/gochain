package chain

import (
	"../utils"
	"fmt"
  "math"
	"strconv"
	"time"
)

const (
	MAX_UINT = ^uint(0)
  HEADER_LOG_PREFIX = "header"
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
  log:= utils.GetLogger(HEADER_LOG_PREFIX)
	if hdr.Height != prevHdr.Height+1 {
    log.Println("Wrong height")
		return false
	}
	if hdr.PrevHash != prevHdr.BlockHash {
    log.Println("PrevHash mismatch")
		return false
	}
	if prevHash := prevHdr.ComputeHash(); hdr.PrevHash != prevHash {
    log.Println("Hash mismatch")
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

func (hdr *BlockHeader) MineNext(difficulty uint, miner Miner) bool {
  logger := utils.GetLogger(HEADER_LOG_PREFIX)
	now := time.Now().UnixNano()
	if miner.MineNext(hdr, difficulty) {
		hdr.Nonce = miner.GetResult()
    hdr.BlockHash = hdr.ComputeHash()
		time_elapsed := float64(time.Now().UnixNano()-now) / math.Pow10(6)
		logger.Printf("nonce: %d", hdr.Nonce)
		logger.Printf("Hashrate: %f", 1000*float64(miner.GetHashProcessed())/time_elapsed)
		logger.Printf("Time elapsed: %f ms", time_elapsed)
	}
	return true
}

func (hdr *BlockHeader) Print() string {
	dataHash := utils.ComputeHashEncoded(hdr.Data)
	return fmt.Sprintf("Height: %d, Timestamp: %s,\nBlockhash: %s, PrevBlockHash: %s\nData hash: %s\nNonce: %d",
		hdr.Height, hdr.Timestamp, hdr.BlockHash, hdr.PrevHash, dataHash, hdr.Nonce)
}
