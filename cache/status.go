package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/carameleon/go-redis-test/client"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	_ "github.com/tendermint/tendermint/rpc/core/types"
)

const (
	displaySize         int64  = 5
	redisQueueSizeLimit int64  = 10
	redisKeyStatus      string = "status"
	redisKeyBlockInfo   string = "block_info"
)

// PushBlockInfo push the latest block
func PushBlockInfo(c context.Context) {
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

		lpush := cli.Rd.LPush(c, redisKeyBlockInfo, data)
		size, err := lpush.Result()
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second)
			continue
		}

		fmt.Println("====================================================================================================")
		fmt.Println("LPush.Size", size)
		fmt.Printf("%-10s : %d\n", "Height", rb.Block.Header.Height)
		// fmt.Printf("%-10s : %s\n", "BlockHash", rb.BlockMeta.BlockID.Hash.String())
		// fmt.Printf("%-10s : %s\n", "Proposer", rb.Block.Header.ProposerAddress.String())
		// fmt.Printf("%-10s : %d\n", "TotalTx", rb.Block.Header.TotalTxs)
		// fmt.Printf("%-10s : %d\n", "Num of Tx", rb.Block.Header.NumTxs)
		// fmt.Printf("%-10s : %s\n", "Time", rb.BlockMeta.Header.Time)

		prevBlockHeight = rb.Block.Header.Height

		if size > 2*redisQueueSizeLimit {
			Trim(c)
		}

		time.Sleep(5 * time.Second)
	}
}

// PopBlockInfo returns list of block info from redis queue
func PopBlockInfo(c context.Context) ([]*BlockInfo, error) {
	cli, ok := IsClient(c.Value("cli"))
	if !ok {
		return nil, errors.Errorf("Context value is not client")
	}
	for cli.Rd.Exists(c, redisKeyBlockInfo).Val() == 0 {
		time.Sleep(time.Second)
		return nil, errors.Errorf("Redis started up just before, retry after 1 sec")
	}
	m := min(displaySize, cli.Rd.LLen(c, redisKeyBlockInfo).Val())
	result := make([]*BlockInfo, m)
	var i int64
	for i = 0; i < m; i++ {
		b, err := cli.Rd.LIndex(c, redisKeyBlockInfo, i).Bytes()
		if err != nil {
			if err == redis.Nil {
				// fmt.Println("pop : data no longer exist")
				// should return if objects are exist.
				return result, nil
			}
			return nil, err
		}
		rb := new(BlockInfo)
		if err = json.Unmarshal(b, rb); err != nil {
			fmt.Println(err)
			continue
		}

		result[i] = rb

		// fmt.Printf("%30s %s: %d\n", "", "Height", rb.Height)
		// fmt.Printf("%30s %s: %d\n", "", "Tx num", rb.NumTxs)
	}
	return result, nil
}

// FlushAll flush redis buffer
func FlushAll(c context.Context) error {
	cli, ok := IsClient(c.Value("cli"))
	if !ok {
		return errors.Errorf("Context value is not client")
	}
	s := cli.Rd.FlushAll(c)
	fmt.Println(s.Result())
	return nil
}

// Trim trim the redis message queue larger than redisQueueSizeLimit
func Trim(c context.Context) {
	cli, ok := IsClient(c.Value("cli"))
	if !ok {
		return
	}

	status := cli.Rd.LTrim(c, redisKeyBlockInfo, 0, redisQueueSizeLimit)
	r, err := status.Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("====================================================================================================")
	fmt.Println("redis queue is trimed :", r)
	fmt.Println("====================================================================================================")
}

// IsClient for type assertion client
func IsClient(i interface{}) (*client.Client, bool) {
	switch v := i.(type) {
	case *client.Client:
		return v, true
	default:
		return nil, false
	}
}

func min(a, b int64) int64 {
	if a >= b {
		return b
	}
	return a
}
