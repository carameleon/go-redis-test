package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const expiredTime = 5 * time.Minute

// SetExpiredCache push the latest block hash that expired specific time duration.
func SetExpiredCache(c context.Context) {
	cli, ok := IsClient(c.Value("cli"))
	if !ok {
		fmt.Println("Context value is not client")
	}
	var prevBlockHeight int64
	for {
		rb, err := cli.RPC.Block(nil)
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second)
			continue
		}

		if prevBlockHeight == rb.Block.Header.Height {
			rb = nil
			fmt.Println("<!> waiting... new block")
			time.Sleep(time.Second)
			continue
		}

		bi := BlockInfo{
			BlockHash: rb.BlockMeta.BlockID.Hash.String(),
			Proposer:  rb.Block.Header.ProposerAddress.String(),
			Height:    rb.Block.Header.Height,
			TotalTxs:  rb.Block.Header.TotalTxs,
			NumTxs:    rb.Block.Header.NumTxs,
			Time:      rb.BlockMeta.Header.Time,
		}

		data, err := json.Marshal(bi)
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second)
			continue
		}

		// var tx GeneralTx

		cli.Rd.Set(c, bi.BlockHash, data, expiredTime)

		prevBlockHeight = rb.Block.Header.Height

		fmt.Println("redis all keys:", cli.Rd.Keys(c, "*"))
		fmt.Println("db size:", cli.Rd.DBSize(c))
		time.Sleep(5 * time.Second)
	}
}
