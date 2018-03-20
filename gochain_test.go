package main

import (
	"./chain"
	"./utils"
	"log"
	"testing"
  "strconv"
  "encoding/base64"
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
	log.Println("Chain test")
	var newChain chain.Chain
	var currentHdr *chain.BlockHeader = nil
	for i := 0; i < 20; i++ {
		b, _ := utils.ReadRandom(1024 * 1024)
		_, nextHdr := chain.NewBlockHeader(currentHdr, b)
		currentHdr = &nextHdr
		newChain = append(newChain, *currentHdr)
	}
  for i := range newChain {
    log.Println(newChain[i].Print())
  }
	log.Println(newChain.Verify())
	newChain[1].Data[2] = 0
	log.Println(newChain.Verify())
}

const MineDifficulty = 24

func TestMine(t *testing.T) {
	b, _ := utils.ReadRandom(50)
	_, hdr := chain.NewBlockHeader(nil, b)
	b, _ = utils.ReadRandom(50)
	_, nextHdr := chain.NewBlockHeader(&hdr, b)
  nextHdr.MineNext(MineDifficulty)
  if !nextHdr.Verify(&hdr) {
    t.Fail()
  }
  log.Println(nextHdr.Print())
  hash, err := base64.StdEncoding.DecodeString(nextHdr.BlockHash)
  if err!= nil {
    t.Errorf("Error decoding: %s", err)
  }
  log.Println(hash)
}

func TestDifficultyCheck (t *testing.T) {
  MyT:= (*MyT)(t)
  hash := make([]byte, 50)
  for i:= range hash {
    hash[i] = 0xff
  }
  i, err := strconv.ParseInt("00011111", 2, 8)
  if err!= nil {
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
