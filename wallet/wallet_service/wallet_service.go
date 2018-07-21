package wallet

import (
	"pmdgo/conf"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/shopspring/decimal"
	"context"
	"pmdgo/services/wallet/tokens"
	"pmdgo/services/wallet/connector/go_ethereum_conn"
	"pmdgo/services/wallet/connector/ethrpc"
	"github.com/garyburd/redigo/redis"
	"fmt"
	"io/ioutil"
	"strings"
	"math/big"
)



type WalletService struct {
	goethClient  *go_ethereum_conn.EthConnector
	ethRpc *ethrpc.EthRPC
	redisPool *redis.Pool
	adminAddr string
	adminKeyFile string
	adminPwd string
}

func NewWalletService(client *go_ethereum_conn.EthConnector, ethRpc *ethrpc.EthRPC,redis_pool *redis.Pool,adminAddr, adminKeyFile,adminPwd string) *WalletService {
	return &WalletService{client,ethRpc,redis_pool, adminAddr,adminKeyFile,adminPwd}
}

func (s *WalletService) GetBalanceErc20(tokenCfg *conf.Token, addr string) (b string,err error) {
	tokenAddress := common.HexToAddress(tokenCfg.TokenAddr)
	c,err := tokens.NewTokenCaller(tokenAddress,s.goethClient.Client())
	if err != nil {
		return
	}

	accountAddress := common.HexToAddress(addr)
	balance,err := c.BalanceOf(&bind.CallOpts{Pending: false,Context:context.Background()},accountAddress)
	if err != nil {
		return
	}

	tokenDecimal := decimal.New(1,tokenCfg.DecimalNum)
	balance_decimal := decimal.NewFromBigInt(balance,0)
	bb := balance_decimal.Div(tokenDecimal)
	return bb.String(),nil
}


func (s *WalletService) GetEthBalance(addr string) (b string, err error) {
	balance, err := s.ethRpc.EthGetBalanceV2(addr, "latest")
	if err != nil {
		return
	}

	d1 := decimal.NewFromBigInt(&balance,1)
	eth_decimal := decimal.New(1,18)
	return d1.Div(eth_decimal).String(),nil
}



func (s *WalletService) AdminTransferToken(to,tokenAddr string,amount float64) (string,error) {
	conn := s.redisPool.Get()
	defer conn.Close()
	isNewNonce := false
	nonce,err := getTransNonce(conn)
	if err != nil && err != redis.ErrNil {
		return "",err
	}

	fromAddr := common.HexToAddress(s.adminAddr)
	if err == redis.ErrNil {
		err = nil
		nonce, err = s.goethClient.Client().PendingNonceAt(s.goethClient.Ctx(), fromAddr)
		if err != nil || nonce < 1 {
			nonce, err = s.goethClient.Client().NonceAt(s.goethClient.Ctx(), fromAddr, nil)
			if err != nil {
				return "",err
			}
		}

		isNewNonce = true
		setTransNonce(nonce, conn)
	}

	if !isNewNonce {
		nonce = nonce + 1
	}

	i,err := ioutil.ReadFile(s.adminKeyFile)
	if err != nil {
		return "",err
	}

	auth,err := bind.NewTransactor(strings.NewReader(string(i)), s.adminPwd)
	if err != nil {
		return "",err
	}

	auth.GasLimit = 100000
	auth.Nonce = big.NewInt(int64(nonce))
	tokenAddress := common.HexToAddress(tokenAddr)
	tk,err := tokens.NewToken(tokenAddress,s.goethClient.Client())
	if err != nil {
		return "",err
	}

	tt,err := tk.Transfer(auth,common.HexToAddress(to), FloatToBigInt(amount))
	return tt.Hash().String(),err
}


func transNonceKey() string {
	return fmt.Sprintf("trans:nonce")
}

func getTransNonce(conn redis.Conn) (uint64, error) {
	return redis.Uint64(conn.Do("GET", transNonceKey()))
}

func setTransNonce(nonce uint64,conn redis.Conn) (int, error) {
	return redis.Int(conn.Do("SET", transNonceKey(), nonce))
}