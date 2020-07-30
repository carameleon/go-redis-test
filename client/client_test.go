package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/carameleon/go-redis-test/config"
	"github.com/go-redis/redis/v8"
)

//run first your redis-server
func TestPing(t *testing.T) {

	ctx := context.Background()

	cfg := config.ParseConfig()

	cli, err := NewClient(cfg)
	if err != nil {
		t.Log(err)
	}

	pong, err := cli.Rd.Ping(ctx).Result()
	if err != nil {
		t.Log(err)
	}
	t.Log(pong)

}

func TestSetData(t *testing.T) {
	ctx := context.Background()
	cfg := config.ParseConfig()

	cli, err := NewClient(cfg)
	if err != nil {
		t.Log(err)
	}

	err = cli.Rd.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := cli.Rd.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("rd.Get(key) :", val)

	val2, err := cli.Rd.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}

func TestLcdPing(t *testing.T) {
	cfg := config.ParseConfig()

	cli, err := NewClient(cfg)
	if err != nil {
		t.Log(err)
	}

	resp, err := cli.Lcd.R().Get("/node_info")
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Log(string(resp.Body()))
	// data, err := resp.Body()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// t.Log(string(data))
}
