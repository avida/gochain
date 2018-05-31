package db

import (
	"../utils"
	"bytes"
	"log"
	"os"
	"testing"
	"time"
)

const TEST_LOG_PREFIX = "test"

var logger *log.Logger

func TestMain(m *testing.M) {
	utils.SetupLoggers()
	logger = utils.GetLogger(TEST_LOG_PREFIX)
	utils.ReadConf("../test")
	os.Exit(m.Run())
}

func TestSimpleDBWrite(t *testing.T) {
	logger.Println("TestDb")
  err := Connect()
  if err != nil {
    t.Errorf("Error connecting to db: %v", err)
  }

	randomBytes, _ := utils.ReadRandom(2 * 10)
	chain := &BlockHeader{
		Height:    1,
		Nonce:     1,
		Timestamp: time.Now().UTC().Format(time.ANSIC),
		BlockHash: "AAAAV2Is6HFTKHPEQklPIDzyt/SVRYOxgTSj0aJXFFE=",
		Data:      randomBytes,
	}

	if err := SaveBlock(chain); err != nil {
		t.Fatalf("Error saving db: %v", err)
	}
	err, restored := LoadBlock(1)
	if err != nil {
		t.Fatalf("Error looking up record: %v", err)
	}

	if restored.Timestamp != chain.Timestamp {
		t.Errorf("Error matching data w/ timezones")
	}

	if restored.BlockHash != chain.BlockHash {
		t.Errorf("Error matching data w/ hashes")
	}

	if !bytes.Equal(restored.Data, chain.Data) {
		t.Fatalf("Error matching data")
	}

	if err := DeleteBlock(restored.Height); err != nil {
		t.Fatalf("Error deleting record %v", err)
	}

	err, _ = LoadBlock(1)
	if err == nil {
		t.Fatalf("Record havent been deleted")
	}
}
