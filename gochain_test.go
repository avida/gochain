package main

import (
	"./chain"
	"./utils"
	"log"
	"testing"
)

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
	log.Println("Chain")
	var newChain chain.Chain
	var currentHdr *chain.BlockHeader = nil
	for i := 0; i < 1000; i++ {
		b, _ := utils.ReadRandom(100)
		_, nextHdr := chain.NewBlockHeader(currentHdr, b)
		currentHdr = &nextHdr
		newChain = append(newChain, *currentHdr)
	}
	log.Println(newChain.Verify())
	newChain[1].Data[2] = 0
	log.Println(newChain.Verify())
}
