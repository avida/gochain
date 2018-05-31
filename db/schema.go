package db

//go:generate reform

//reform:chain
type BlockHeader struct {
	Height    int    `reform:"height,pk"`
	Nonce     uint   `reform:"nonce"`
	Timestamp string `reform:"timestamp"`
	BlockHash string `reform:"block_hash"`
	PrevHash  string `reform:"prev_hash"`
	Data      []byte `reform:"data"`
}
