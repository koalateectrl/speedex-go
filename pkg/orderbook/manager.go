package orderbook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"

	"github.com/sandymule/speedex-go/pkg/assets"
	"github.com/sandymule/speedex-go/pkg/dmdutils"
)

type Manager map[assets.AssetPair]Orderbook

func (mgr Manager) SpyDmdQuery(prices map[assets.Asset]float64, smoothMult uint8) dmdutils.SpyDmd {
	sd := make(dmdutils.SpyDmd)

	for assetPair, ob := range mgr {
		sellPrice := prices[assetPair.Sell]
		buyPrice := prices[assetPair.Buy]
		tradeAmt := ob.GetAmt(sellPrice, buyPrice, smoothMult)
		sd.AddSpyDmdPair(assetPair, tradeAmt)
	}

	return sd
}

func (mgr *Manager) AddFromJson(path string) {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println("Successfully Opened txs.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var txs Txs
	json.Unmarshal(byteValue, &txs)

	for i := 0; i < len(txs.Txs); i++ {
		fmt.Println("CALLED")
		mgr.AddOneTx(txs.Txs[i])
	}
}

func (mgr *Manager) AddOneTx(tx Tx) {
	var ap assets.AssetPair
	ap.Buy = assets.Asset(tx.BuyType)
	ap.Sell = assets.Asset(tx.SellType)

	// if AssetPair not in Manager, start with 0 price entry
	if _, ok := (*mgr)[ap]; !ok {
		var ob Orderbook
		// start with 0 price entry
		ob = append(ob, PriceCompStats{SellPrice: 0, CumForSale: 0, CumForSaleTimesPrice: 0})
		(*mgr)[ap] = ob
	}

	var pcs PriceCompStats
	pcs.SellPrice = tx.SellLimitPrice
	pcs.Txid = tx.Txid

	currOb := (*mgr)[ap]
	pos := sort.Search(len(currOb), func(i int) bool { return currOb[i].SellPrice >= tx.SellLimitPrice })
	pcs.CumForSale = currOb[pos-1].CumForSale + tx.SellAmount
	pcs.CumForSaleTimesPrice = currOb[pos-1].CumForSaleTimesPrice + tx.SellAmount*tx.SellLimitPrice

	// insert transaction into orderbook
	newob := make(Orderbook, len(currOb)+1)
	at := copy(newob, currOb[:pos])
	newob[pos] = pcs
	at++
	copy(newob[at:], currOb[pos:])

	// Update cumulative values for all entries with higher sell limit price
	for j := pos + 1; j < len(newob); j++ {
		fmt.Println("RAN")
		newob[j].CumForSale += tx.SellAmount
		newob[j].CumForSaleTimesPrice += tx.SellAmount * tx.SellLimitPrice
	}

	(*mgr)[ap] = newob
}

func (mgr Manager) GetExecTxs(prices map[assets.Asset]float64, smoothMult uint8) map[assets.AssetPair]ExecTxBook {
	execTxs := make(map[assets.AssetPair]ExecTxBook)
	for assetPair, ob := range mgr {
		execTxB := make(ExecTxBook, 0)
		execPrice := prices[assetPair.Sell] / prices[assetPair.Buy]
		fullPrice := ApplySmoothMult(execPrice, smoothMult)
		canSell := true

		for idx, tx := range ob {
			if idx == 0 {
				continue
			}

			var execTxStats ExecTxStats
			execTxStats.Txid = tx.Txid
			execTxStats.SellPrice = tx.SellPrice

			if tx.SellPrice > execPrice {
				canSell = false
			}

			if canSell {
				if tx.SellPrice < fullPrice {
					execTxStats.AmtSold = tx.CumForSale - ob[idx-1].CumForSale
					execTxStats.AmtRem = 0
				} else {
					execTxStats.AmtSold = (execPrice - tx.SellPrice) / (execPrice * math.Pow(2, -float64(smoothMult))) * (tx.CumForSale - ob[idx-1].CumForSale)
					execTxStats.AmtRem = tx.CumForSale - ob[idx-1].CumForSale - execTxStats.AmtSold
				}

			} else {
				execTxStats.AmtSold = 0
				execTxStats.AmtRem = tx.CumForSale - ob[idx-1].CumForSale
			}
			execTxB = append(execTxB, execTxStats)
		}

		execTxs[assetPair] = execTxB

	}

	return execTxs
}
