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

	do = make(chan int)
	done = make(chan int)

	cfg := config.ParseConfig()

	cli, err := client.NewClient(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "cli", cli)

	time.Sleep(time.Second)

	go cache.SetStatus(ctx, do)
	go cache.GetStatus(ctx, 50)
	go cache.Timer(do, 5)

	<-done
}
