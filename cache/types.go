package cache

import (
	"encoding/json"
	"time"
)

// BlockInfo has block information
type BlockInfo struct {
	BlockHash string    `json:"block_hash" sql:",unique"`
	Height    int64     `json:"height"`
	Proposer  string    `json:"proposer"`
	TotalTxs  int64     `json:"total_txs" sql:"default:0"`
	NumTxs    int64     `json:"num_txs" sql:"default:0"`
	Time      time.Time `json:"time"`
}

// GeneralTx is a struct for general tx
type GeneralTx struct {
	Height    string          `json:"height"`
	TxHash    string          `json:"txhash"`
	Data      string          `json:"data"`
	RawLog    string          `json:"raw_log"`
	Logs      json.RawMessage `json:"logs"`
	GasWanted string          `json:"gas_wanted"`
	GasUsed   string          `json:"gas_used"`
	Tags      json.RawMessage `json:"tags"`
	Tx        json.RawMessage `json:"tx"`
	Timestamp time.Time       `json:"timestamp"`
}
