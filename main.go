package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/carameleon/go-redis-test/cache"
	"github.com/carameleon/go-redis-test/client"
	"github.com/carameleon/go-redis-test/config"
)

var do chan int
var done chan int

func main() {

	done = make(chan int)

	cfg := config.ParseConfig()

	cli, err := client.NewClient(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "cli", cli)

	cache.FlushAll(ctx)
	go cache.SetExpiredCache(ctx)
	go cache.PushBlockInfo(ctx)

	for i := 0; i < 1; i++ {
		go func() {
			for {
				// fmt.Println("====================================================================================================")
				result, err := cache.PopBlockInfo(ctx)
				if err != nil {
					fmt.Println(err)
				}
				resultPrint(result)
				time.Sleep(10 * time.Millisecond)
				// time.Sleep(1 * time.Second)
				// repond(rw, result)
			}
		}()
	}

	<-done
}

func resultPrint(result []*cache.BlockInfo) {
	return
	str := "result:"
	for _, b := range result {
		str = fmt.Sprintf("%s %d", str, b.Height)
	}
	fmt.Println(str)
}
