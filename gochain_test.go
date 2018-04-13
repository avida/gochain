package main

import (
	"./chain"
	"./utils"
	"encoding/base64"
	"log"
	"strconv"
	"testing"
	//"github.com/davecgh/go-spew/spew"
)

type MyT testing.T

func (t *MyT) checkTrue(condition bool, errorMsg string) {
	if !condition {
		t.Error(errorMsg)
	}
}

func TestMain(t *testing.T) {
	t.Log("test")
	log.Println(utils.ComputeHash([]byte(FirstData)))
	log.Println("tests")
	b, _ := utils.ReadRandom(5)
	log.Println(b)
}

func TestCreateHeader(t *testing.T) {
	b, _ := utils.ReadRandom(50)
	_, hdr := chain.NewBlockHeader(nil, b)
	b, _ = utils.ReadRandom(50)
	_, nextHdr := chain.NewBlockHeader(&hdr, b)
	_, thirdHdr := chain.NewBlockHeader(&nextHdr, b)
	log.Println("Test header")
	log.Printf("header is %v .", hdr)
	log.Printf("header is %v .", nextHdr)
	log.Printf("header is %v .", thirdHdr)
}

func TestChain(t *testing.T) {
	MyT := (*MyT)(t)
	log.Println("Chain test")
  const Blocks = 20
  const Difficulty = 16
	var ledger chain.Chain
	for i := 0; i < Blocks; i++ {
    block := ledger.AddBlock()
    var miner chain.MultiThreadMiner
    block.MineNext(Difficulty, &miner)
	}
  MyT.checkTrue(ledger.Verify(), "Check ledger")
  log.Println(len(ledger))
  for block:= range ledger {
    hash, err := base64.StdEncoding.DecodeString(ledger[block].BlockHash)
    MyT.checkTrue(err == nil, "Hash decoded ok")
    MyT.checkTrue(chain.CheckHashOk(hash, Difficulty), "Check block number " + string(block))
  }
	ledger[7].Data[2] = 0
	MyT.checkTrue(ledger.Verify() == false, "Check modified ledger fail")
}

const MineDifficulty = 24

func MakeBlock(prevBlock *chain.BlockHeader, data_size uint ) (err error, block *chain.BlockHeader) {
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
  _, firstBlock:= MakeBlock(nil, 50)
  _, secondBlock:= MakeBlock(firstBlock, 50)
	var miner chain.MultiThreadMiner
	secondBlock.MineNext(MineDifficulty, &miner)
	if !secondBlock.Verify(firstBlock) {
		t.Fatal("Block verification failed")
	}
	hash, err := base64.StdEncoding.DecodeString(secondBlock.BlockHash)
	if err != nil || !chain.CheckHashOk(hash, MineDifficulty) {
		t.Fatal("Difficulty doesnt match")
	}
	log.Println(secondBlock.Print())
	log.Println(hash)
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
  _, firstBlock:= MakeBlock(nil, 50)
  _, secondBlock:= MakeBlock(firstBlock, 50)
	var miner chain.MultiThreadRangeMiner
	secondBlock.MineNext(MineDifficulty, &miner)
  MyT.checkTrue(secondBlock.Verify(firstBlock),
  "Check block validity")
	hash, _:= base64.StdEncoding.DecodeString(secondBlock.BlockHash)
  MyT.checkTrue(chain.CheckHashOk(hash, MineDifficulty),
  "Chechk difficulty matches")
  log.Println(hash)
}
