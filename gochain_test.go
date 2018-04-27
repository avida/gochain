package main

import (
	"./chain"
	"./utils"
	"encoding/base64"
	"log"
	"os"
	"strconv"
	"testing"
)

const (
	TEST_LOG_PREFIX = "test"
	MINE_DIFFICULTY = 24
)

var logger *log.Logger

type MyT testing.T

func (t *MyT) checkTrue(condition bool, errorMsg string) {
	if !condition {
		t.Error(errorMsg)
	}
}

func TestMain(m *testing.M) {
	SetupLoggers()
	logger = utils.GetLogger(TEST_LOG_PREFIX)
	os.Exit(m.Run())
}

func TestHash(t *testing.T) {
	t.Log("test")
	logger.Println(utils.ComputeHash([]byte(FirstData)))
	logger.Println("tests")
	b, _ := utils.ReadRandom(5)
	logger.Println(b)
}

func TestCreateHeader(t *testing.T) {
	b, _ := utils.ReadRandom(50)
	_, hdr := chain.NewBlockHeader(nil, b)
	b, _ = utils.ReadRandom(50)
	_, nextHdr := chain.NewBlockHeader(&hdr, b)
	_, thirdHdr := chain.NewBlockHeader(&nextHdr, b)
	logger.Println("Test header")
	logger.Printf("header is %v .", hdr)
	logger.Printf("header is %v .", nextHdr)
	logger.Printf("header is %v .", thirdHdr)
}

func TestChain(t *testing.T) {
	MyT := (*MyT)(t)
	logger.Println("Chain test")
	const Blocks = 30
	const Difficulty = 10
	var ledger chain.Chain
	for i := 0; i < Blocks; i++ {
		block := ledger.AddBlock()
		miner := chain.MakeMultiThreadMiner(8)
		block.MineNext(Difficulty, miner)
	}
	MyT.checkTrue(ledger.Verify(), "Check ledger")
	for block := range ledger {
		hash, err := base64.StdEncoding.DecodeString(ledger[block].BlockHash)
		MyT.checkTrue(err == nil, "Hash decoded ok")
		MyT.checkTrue(chain.CheckHashOk(hash, Difficulty), "Check block number "+string(block))
	}
	ledger[7].Data[2] = 0
	MyT.checkTrue(ledger.Verify() == false, "Check modified ledger fail")
}

func MakeBlock(prevBlock *chain.BlockHeader, data_size uint) (err error, block *chain.BlockHeader) {
	b, e := utils.ReadRandom(50)
	if e != nil {
		return e, nil
	}
	e, hdr := chain.NewBlockHeader(prevBlock, b)
	if e != nil {
		return e, nil
	}
	return nil, &hdr
}

func TestMine(t *testing.T) {
	_, firstBlock := MakeBlock(nil, 50)
	_, secondBlock := MakeBlock(firstBlock, 50)
	miner := chain.MakeMultiThreadMiner(8)
	secondBlock.MineNext(MINE_DIFFICULTY, miner)
	if !secondBlock.Verify(firstBlock) {
		t.Fatal("Block verification failed")
	}
	hash, err := base64.StdEncoding.DecodeString(secondBlock.BlockHash)
	if err != nil || !chain.CheckHashOk(hash, MINE_DIFFICULTY) {
		t.Fatal("Difficulty doesnt match")
	}
	logger.Println(secondBlock.Print())
	logger.Println(hash)
}

func TestDifficultyCheck(t *testing.T) {
	MyT := (*MyT)(t)
	hash := make([]byte, 50)
	for i := range hash {
		hash[i] = 0xff
	}
	i, err := strconv.ParseInt("00011111", 2, 8)
	if err != nil {
		t.Errorf("Error parsing: %s", err)
	}
	hash[0] = 0
	hash[1] = 0
	hash[2] = byte(i)
	MyT.checkTrue(chain.CheckHashOk(hash, 19),
		"Difficulty check failed")
	MyT.checkTrue(!chain.CheckHashOk(hash, 20),
		"Difficulty check failed")
}

func TestRangeMiner(t *testing.T) {
	MyT := (*MyT)(t)
	_, firstBlock := MakeBlock(nil, 50)
	_, secondBlock := MakeBlock(firstBlock, 50)
	miner := chain.MakeRangeMiner(8)
	secondBlock.MineNext(MINE_DIFFICULTY, miner)
	MyT.checkTrue(secondBlock.Verify(firstBlock),
		"Check block validity")
	hash, _ := base64.StdEncoding.DecodeString(secondBlock.BlockHash)
	MyT.checkTrue(chain.CheckHashOk(hash, MINE_DIFFICULTY),
		"Check difficulty matches")
	logger.Println(hash)
}
