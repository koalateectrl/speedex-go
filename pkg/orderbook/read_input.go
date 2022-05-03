package orderbook

type Txs struct {
	Txs []Tx `json:"txs"`
}

type Tx struct {
	Txid           uint64  `json:"txid"`
	BuyType        string  `json:"buytype"`
	SellType       string  `json:"selltype"`
	SellAmount     float64 `json:"sellamount"`
	SellLimitPrice float64 `json:"selllimitprice"`
}
