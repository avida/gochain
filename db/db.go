package db

import (
	"../utils"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
	"strings"
)

const (
	DB_LOG_PREFIX = "db"
)

var db *reform.DB

func connStr() string {
	port, user, pass := utils.DBConfig()
	return fmt.Sprintf("postgresql://localhost:%d/%v?user=%v&password=%v&sslmode=disable", port, user, user, pass)
}

func Connect() error {
	logger := utils.GetLogger(DB_LOG_PREFIX)
	if db != nil {
		return errors.New("Already connected")
	}
	conn, err := sql.Open("postgres", connStr())
	if err == nil {
		err = conn.Ping()
	}
	if err != nil {
		return err
	}
	db = reform.NewDB(conn, postgresql.Dialect, reform.NewPrintfLogger(logger.Printf))
	if db == nil {
		return errors.New("Error creating reform db")
	}

	return nil
}

func SaveBlock(block *BlockHeader) error {
	return db.Save(block)
}

func LoadBlock(height int) (error, *BlockHeader) {
	restored, err := db.FindByPrimaryKeyFrom(BlockHeaderTable, height)
	if err == nil {
		// Strings that are shorter than 40 symbols (defined in schema.sql) are padded with the whitespaces.
		restored.(*BlockHeader).BlockHash = strings.Trim(restored.(*BlockHeader).BlockHash, " ")
		restored.(*BlockHeader).PrevHash = strings.Trim(restored.(*BlockHeader).PrevHash, " ")
		restored.(*BlockHeader).Timestamp = strings.Trim(restored.(*BlockHeader).Timestamp, " ")
		return nil, restored.(*BlockHeader)
	}
	return err, nil
}

func DeleteBlock(height int) error {
	block := BlockHeaderTable.NewRecord()
	block.(*BlockHeader).Height = height
	return db.Delete(block)
}
