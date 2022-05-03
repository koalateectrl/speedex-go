package main

import (
	"fmt"

	"github.com/sandymule/speedex-go/pkg/assets"
	"github.com/sandymule/speedex-go/pkg/orderbook"
	"github.com/sandymule/speedex-go/pkg/tatonnement"
)

func main() {

	var oracle tatonnement.Oracle

	cp := tatonnement.CtrlParams{MSmoothMult: 15, MMaxRnds: 1000,
		MStepUp: 40, MStepDown: 25, MStepSizeRdx: 5, MStepRdx: 30}
	prices := make(map[assets.Asset]float64)

	/*prices["ETH"] = 4500
	prices["USDT"] = 1*/

	prices["USD"] = 1
	prices["EUR"] = 1.2

	mgr := make(orderbook.Manager)
	mgr.AddFromJson("/Users/samwwong/Desktop/speedex-go/test_cases/txs2.json")

	oracle.ObManager = mgr
	newPrices := oracle.ComputePrices(cp, prices, 5)
	fmt.Println(newPrices)

	execTxs := mgr.GetExecTxs(newPrices, cp.MSmoothMult)
	fmt.Println(execTxs)
}
