package wallet



type TransferToken struct {
	coin string `json:"coin" binding:"required"`
	ToAddr string `json:"to_addr" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}
