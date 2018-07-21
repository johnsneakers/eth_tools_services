package go_ethereum_conn

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"context"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)


type EthConnector struct {
	ctx       context.Context
	ethclient *ethclient.Client
	chainId *big.Int
}

func InitEthConnector(url string, chainId int64) *EthConnector {
	var ethConn EthConnector
	ctx := context.Background()
	c, err := ethrpc.DialContext(ctx, url)
	if err != nil {
		panic(err)
	}

	ethConn.ethclient = ethclient.NewClient(c)
	ethConn.ctx = ctx
	ethConn.chainId = big.NewInt(chainId)
	return &ethConn
}


func (c *EthConnector) Client() *ethclient.Client {
	return c.ethclient
}

func (c *EthConnector) Ctx() context.Context {
	return c.ctx
}


func (c *EthConnector) GetDefaultGasPrice() *big.Int {
	return new(big.Int).SetInt64(1000000000 * 55)
}

func (c *EthConnector) GetDefaultGasLimit() *big.Int {
	return new(big.Int).SetInt64(250000)
}
