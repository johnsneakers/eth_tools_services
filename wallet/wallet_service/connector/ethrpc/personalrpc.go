package ethrpc

import "fmt"

func (rpc *EthRPC) Personal_newAccount(password string) (string, error) {
	var hash string
	fmt.Println("创建===>",rpc.url)
	err := rpc.call("personal_newAccount", &hash, password)
	return hash, err
}