package client

import (
	"fmt"
	"time"

	"github.com/carameleon/go-redis-test/config"
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	trpc "github.com/tendermint/tendermint/rpc/client"
)

// Client implements a wrapper around both a Tendermint RPC client and a
// Cosmos SDK REST client that allows for essential data queries.
type Client struct {
	RPC *trpc.HTTP
	Lcd *resty.Client
	Rd  *redis.Client
}

// NewClient creates a new client with the given config
func NewClient(cfg *config.Config) (*Client, error) {

	rpc := trpc.NewHTTP(cfg.Node.RPC, "/websocket")

	lcd := resty.New().
		SetHostURL(cfg.Node.LCD).
		SetTimeout(time.Duration(10 * time.Second))

	redisAddr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	rd := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &Client{rpc, lcd, rd}, nil
}
