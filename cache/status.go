package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/carameleon/go-redis-test/client"
	"github.com/go-redis/redis/v8"
)

func SetStatus(c context.Context, do <-chan int) {
	var i int
	var value string
	cli := c.Value("cli").(*client.Client)

	for {
		select {
		case <-do:
			// resp, err := cli.Lcd.R().Get("/blocks/latest")
			resp, err := cli.Lcd.R().Get("/blocks/latest")
			if err != nil {
				fmt.Println(err)
				continue
			}
			data, err := json.Marshal(resp.Body())
			err = cli.Rd.Set(c, "block_latest", data, 0).Err()
			if err != nil {
				panic(err)
			}
			fmt.Println("set now", value)
			i = 0
		default:
			fmt.Println("waiting data...", i)
			time.Sleep(1 * time.Second)
			i++
		}
	}
}

func GetStatus(c context.Context, t time.Duration) {
	cli := c.Value("cli").(*client.Client)
	for {
		val, err := cli.Rd.Get(c, "block_latest").Result()
		if err != nil {

			if err == redis.Nil {
				fmt.Println("key2 does not exist")
				continue
			}
			fmt.Println(err)
		}
		fmt.Println("Get(key) :", val[:100])
		time.Sleep(t * time.Millisecond)
	}
}

func Timer(do chan<- int, t time.Duration) {
	for {
		do <- 1
		time.Sleep(t * time.Second)
	}
}
