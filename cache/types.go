package cache

import "time"

// BlockInfo has block information
type BlockInfo struct {
	BlockHash string    `json:"block_hash" sql:",unique"`
	Height    int64     `json:"height"`
	Proposer  string    `json:"proposer"`
	TotalTxs  int64     `json:"total_txs" sql:"default:0"`
	NumTxs    int64     `json:"num_txs" sql:"default:0"`
	Time      time.Time `json:"time"`
}
