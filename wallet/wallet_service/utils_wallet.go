package wallet

import (
	"strings"
	"pmdgo/conf"
	"math/big"
)

func (s *service) GetCoinCfg(coin string) (isSupply bool,cfg *conf.Token){
	for _,v := range s.walletCfg.Token {
		coin = strings.ToLower(coin)
		if v.Name == coin {
			return true,v
		}
	}

	return false,nil
}

func FloatToBigInt(val float64) *big.Int {
	bigval := new(big.Float)
	bigval.SetFloat64(val)

	coin := new(big.Float)
	coin.SetInt(big.NewInt(1000000000000000000))

	bigval.Mul(bigval, coin)

	result := new(big.Int)
	bigval.Int(result)

	return result
}