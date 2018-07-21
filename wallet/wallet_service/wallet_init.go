package wallet

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/garyburd/redigo/redis"
	"pmdgo/services"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"fmt"
	"pmdgo/common/errors"
	"net/http"
	"pmdgo/common"
	"github.com/gin-contrib/cors"
	"pmdgo/common/lib"
	cconf "pmdgo/conf"
	"pmdgo/services/wallet/connector/go_ethereum_conn"
	"pmdgo/services/wallet/connector/ethrpc"
)

type service struct {
	conf        *cconf.ServiceConf
	db          *sql.DB
	rp          *redis.Pool
	resp        *common.Resp
	walletCfg   *cconf.AdminWalletConf
	walletService *WalletService
}

const (
	VERSION = "0.0.1"
)

func NewService() services.Service {
	return new(service)
}

func (s *service) Name() string {
	return "crawl"
}

func (s *service) Version() string {
	return VERSION
}

func (s *service) Config(conf *cconf.ServiceConf) error {
	s.conf = conf
	s.walletCfg = cconf.LoadWalletConf()
	s.rp = lib.NewPool(s.conf.Redis.Host,500, s.conf.Redis.Port,10)
	s.EthConnectInit()
	go s.HttpConfig()
	return nil
}

func (s *service) EthConnectInit () {
	url := s.walletCfg.OutUrl
	ethRpc := ethrpc.NewEthRPC(url)
	goEthConn := go_ethereum_conn.InitEthConnector(url,s.walletCfg.Network)
	s.walletService = NewWalletService(goEthConn,ethRpc,s.rp,s.walletCfg.Address,s.walletCfg.KeyFile,s.walletCfg.Pwd)
}

func (s *service) HttpConfig() {
	r := gin.Default()
	r.Use(catchError())
	r.Use(cors.Default())
	r.Use(func(c *gin.Context) { s.resp = common.NewResp(c) })
	r1 := r.Group("/api")
	r1.GET("/balance/:addr", s.GetAllBalance)
	r1.POST("/tranfer/token",s.TransferToken)
	r1.POST("/miner")
	r1.GET("/miner/future")
	addr := fmt.Sprintf("%s:%d", s.conf.Server.Host, s.conf.Server.Port)
	r.Run(addr)
}


func (s *service) Start() {
	color.Green("%s now start!", s.Name())
	select{}
}



func catchError() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				switch err.(type) {
				case *errors.Error:
					e := err.(*errors.Error)
					c.JSON(e.HTTPStatus, e)
					c.Abort()
				default:
					panic(err)
					c.JSON(http.StatusBadRequest, common.CommonErr(-500, fmt.Sprintf("system err:%v",err)))
				}
			}
		}()
		c.Next()
	}
}