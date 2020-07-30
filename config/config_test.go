package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	config := ParseConfig()

	require.NotEmpty(t, config.Node.RPC, "Node LCDEndpoint should not be empty")
	require.NotEmpty(t, config.Node.LCD, "Node LCDEndpoint should not be empty")
	require.NotEmpty(t, config.Redis.Host, "Redis Host is empty")
	require.NotEmpty(t, config.Redis.Port, "Redis Port is empty")
	require.NotEmpty(t, config.Web.Port, "Web Port is empty")

	t.Log("rpc enpoint :", config.Node.RPC)
	t.Log("lcd enpoint :", config.Node.LCD)
	t.Log("redis host :", config.Redis.Host)
	t.Log("redis port :", config.Redis.Port)
	t.Log("web port :", config.Web.Port)
}
