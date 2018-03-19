package chain

import (
	"../utils"
	"fmt"
	// "log"
)

type BlockHeader struct {
	Height              int
	PrevHash, BlockHash string
	Timestamp           string
	Data                []byte
}

func NewBlockHeader(prev *BlockHeader, data []byte) (err error, block BlockHeader) {
	err = nil
	block = BlockHeader{Height: 1}
	if prev != nil {
		block.Height = prev.Height + 1
		block.PrevHash = prev.BlockHash
	}
	block.Timestamp = "Not today"
	block.Data = data
	block.BlockHash = block.ComputeHash()
	return
}
func (hdr *BlockHeader) ComputeHash() string {
	dataHash := utils.ComputeHash(hdr.Data)
	return utils.ComputeHash([]byte(hdr.Timestamp + hdr.PrevHash + dataHash))
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

func (hdr *BlockHeader) Print() string {
	dataHash := utils.ComputeHash(hdr.Data)
	return fmt.Sprintf("Height: %d, Timestamp: %s,\nBlockhash: %s, PrevBlockHash: %s\nData hash: %s",
		hdr.Height, hdr.Timestamp, hdr.BlockHash, hdr.PrevHash, dataHash)
}
