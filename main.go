package main

import (
	"context"
	"fmt"
	"os"

	"github.com/carameleon/go-redis-test/cache"
	"github.com/carameleon/go-redis-test/client"
	"github.com/carameleon/go-redis-test/config"
)

var do chan int
var done chan int

//https://rpc.cluster-galaxynet.iov.one/
func main() {

	do = make(chan int)
	done = make(chan int)

	cfg := config.ParseConfig()

	cli, err := client.NewClient(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	type ctxKey string

	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxKey("cli"), cli)

	go cache.GetStatus(ctx, 50)
	go cache.SetStatus(ctx, do)
	go cache.Timer(do, 5)

	<-done
}
