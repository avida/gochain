package chain

import (
	"../utils"
)

const (
	DATA_SIZE        = 30
	CHAIN_LOG_PREFIX = "chain"
)

type Chain []*BlockHeader

func (chn Chain) Verify() bool {
	logger := utils.GetLogger(CHAIN_LOG_PREFIX)
	for i := range chn {
		invi := len(chn) - i - 1
		if invi == 0 {
			return true
		} else {
			if !chn[invi].Verify(chn[invi-1]) {
				logger.Println("Failed to verify block ", invi)
				logger.Println("Block hash is ", chn[invi-1].BlockHash)
				return false
			}
		}
	}
	return false

}

func (chn *Chain) AddBlock() *BlockHeader {
	var topBlock *BlockHeader
	if len(*chn) == 0 {
		topBlock = nil
	} else {
		topBlock = (*chn)[len(*chn)-1]
	}
	b, _ := utils.ReadRandom(DATA_SIZE)
	_, nextBlock := NewBlockHeader(topBlock, b)
	*chn = append(*chn, &nextBlock)
	return &nextBlock
}
