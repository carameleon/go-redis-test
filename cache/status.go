package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/carameleon/go-redis-test/client"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// SetStatus set the latest block info to redis
func SetStatus(c context.Context, do <-chan int) error {
	cli, ok := IsClient(c.Value("cli"))
	if !ok {
		return errors.Errorf("Context value is not client")
	}
	var i int
	for {
		select {
		case <-do:
			rb, err := cli.RPC.Block(nil)
			if err != nil {
				fmt.Println(err)
				return err
			}

			data, err := json.Marshal(rb)
			if err != nil {
				fmt.Println("set :", err)
				return err
			}

			err = cli.Rd.Set(c, "block_latest", data, 0).Err()
			if err != nil {
				fmt.Println(err)
				return err
			}
			i = 0
		default:
			fmt.Println("next request waiting ...", i)
			time.Sleep(1 * time.Second)
			i++
		}
	}
}

// GetStatus get latest block from redis
func GetStatus(c context.Context, t time.Duration) {
	cli, ok := IsClient(c.Value("cli"))
	if !ok {
		return
	}
	rb := new(ctypes.ResultBlock)
	var p_height int64
	for {
		// val, err := cli.Rd.Get(c, "block_latest").Result()
		b, err := cli.Rd.Get(c, "block_latest").Bytes()
		if err != nil {
			if err == redis.Nil {
				fmt.Println("key does not exist")
				continue
			}
			fmt.Println(err)
		}
		if err = json.Unmarshal(b, rb); err != nil {
			fmt.Println(err)
			continue
		}

		if p_height != rb.Block.Height {
			fmt.Println("Height :", rb.Block.Height)
			fmt.Println("Tx num :", len(rb.Block.Data.Txs))
			for i, tx := range rb.Block.Data.Txs {
				fmt.Println(i, tx.Hash())
			}
		}
		p_height = rb.Block.Height
		time.Sleep(t * time.Millisecond)
	}
}

// Timer for request latest block using rpc-client
func Timer(do chan<- int, t time.Duration) {
	for {
		do <- 1
		time.Sleep(t * time.Second)
	}
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
