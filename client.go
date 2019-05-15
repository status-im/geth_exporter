package main

import (
	"errors"
	"github.com/ethereum/go-ethereum/rpc"
)

type client struct {
	rpcClient *rpc.Client
}

func newClient(ipcPath string) (*client, error) {
	rpcClient, err := rpc.Dial(ipcPath)
	if err != nil {
		return nil, err
	}

	return &client{rpcClient}, nil
}

func (c *client) close() {
	c.rpcClient.Close()
}

func (c *client) metrics() (metrics, error) {
	var res metrics
	return res, c.rpcClient.Call(&res, "debug_metrics", true)
}

func (c *client) currentBlock() (string, error) {
	var res string
	return res, c.rpcClient.Call(&res, "eth_blockNumber")
}

func (c *client) currentBlockMetrics() (metrics, error) {
	// if not syncing we can at least get the current block
	cblock, err := c.currentBlock()
	// and fake syncing metrics from it
	return metrics{
		"currentBlock":  cblock,
		"highestBlock":  cblock,
		"startingBlock": cblock,
		"knownStates":   "0",
		"pulledStates":  "0",
	}, err
}

func (c *client) syncingMetrics() (metrics, error) {
	// using interface{} to check type later
	var res interface{}
	err := c.rpcClient.Call(&res, "eth_syncing")
	// syncing might be turned off, like in statusd
	if err != nil {
		return c.currentBlockMetrics()
	}
	// eth_syncing can return either metrics type, or just a bool
	switch res.(type) {
	case bool: // this bool will always be false
		return c.currentBlockMetrics()
	default:
		// if syncing we need to cast to metrics type
		conv, ok := res.(map[string]interface{})
		if !ok {
			return nil, errors.New("response from eth_syncing unrecognizable")
		}
		return metrics(conv), err
	}
}
