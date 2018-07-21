package wallet

import (
	"github.com/gin-gonic/gin"
	"pmdgo/common/errors"
	"strings"
)

func (s *service) GetAllBalance(c *gin.Context) {
	addr := c.Param("addr")
	coins := c.Query("coins")
	if coins == "" || addr == ""{
		panic(errors.ErrWrongParam)
	}

	coins_arr := strings.Split(coins, ",")
	ret := map[string]string{}
	for _, coin := range coins_arr {
		isSupply,coin_cfg := s.GetCoinCfg(coin)
		if !isSupply {
			continue
		}

		if strings.ToLower(coin) == "eth" {
			balance,err := s.walletService.GetEthBalance(addr)
			if err != nil {
				panic(err)
			}
			ret[coin] = balance
			continue
		}

		balance,err := s.walletService.GetBalanceErc20(coin_cfg, addr)
		if err != nil {
			panic(err)
		}

		ret[coin] = balance
	}
	s.resp.SuccWithData(ret)
	return
}

func (s *service) TransferToken(c *gin.Context) {
	req := &TransferToken{}
	if err := c.ShouldBindJSON(req);err != nil {
		panic(errors.ErrWrongParam)
	}

	isSupply,coin_cfg := s.GetCoinCfg(req.coin)
	if !isSupply {
		panic(errors.ErrCoin)
	}

	hash,err := s.walletService.AdminTransferToken(req.ToAddr,coin_cfg.TokenAddr,req.Amount)
	if err != nil {
		panic(err)
	}


	ret := map[string]string{
		"to":req.ToAddr,
		"tx":hash,
	}

	s.resp.SuccWithDataV2(ret)
	return
}