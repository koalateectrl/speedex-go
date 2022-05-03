package orderbook

import (
	"math"
)

type PriceCompStats struct {
	SellPrice            float64
	CumForSale           float64
	CumForSaleTimesPrice float64
	Txid                 uint64
}

type Orderbook []PriceCompStats

type ExecTxStats struct {
	Txid      uint64
	SellPrice float64
	AmtSold   float64
	AmtRem    float64
}

type ExecTxBook []ExecTxStats

func (ob Orderbook) GetPCS(sellPrice float64, buyPrice float64) PriceCompStats {
	// first entry of the data is the entry where sellPrice is 0
	if len(ob) == 1 {
		return ob[0]
	}

	var start uint8 = 1
	var end uint8 = uint8(len(ob) - 1)

	if ob[end].SellPrice <= sellPrice/buyPrice {
		return ob[end]
	}

	// binary search
	for {
		var mid uint8 = (start + end) / 2
		if start == end {
			return ob[start-1]
		}
		if ob[mid].SellPrice <= sellPrice/buyPrice {
			start = mid + 1
		} else {
			end = mid
		}
	}
}

func ApplySmoothMult(sellPrice float64, smoothMult uint8) float64 {
	if smoothMult == 0 {
		return sellPrice
	}

	return sellPrice - sellPrice*math.Pow(2, -float64(smoothMult))
}

func (ob Orderbook) GetAmt(sellPrice float64, buyPrice float64, smoothMult uint8) float64 {
	/*
		p = (sellPrice / buyPrice)
		An offer executes fully if offer.minPrice < p * (1-2^{smoothMult}),
				 executes not at all if offer.minPrice > p,
		     and executes x-fractionally in the interim
		     for x = (p - offer.minPrice) / (p * 2^{-smoothMult});
		Note that x\in[0,1] and offer execution is a continuous function.
		Considering only offers in the interim gap:
		amount sold = \Sigma_i ((p-offers[i].minPrice) / (p * 2^{-smoothMult}) * offers[i].amount)
					= \Sigma_i (offers[i].amount) / (2^{-smoothMult}) - \Sigma_i (offers[i].amount * offers[i].minPrice) / (p * 2^{-smoothMult})
		amount sold * sellPrice (== amount bought * buyPrice)
			= 2^{smoothMult} * (
					\Sigma_i offers[i].amount * sellPrice
				  - \Sigma_i offers[i].amount * offers[i].minPrice * sellPrice / (sellPrice / buyPrice)
				)
			= 2^{smoothMult} * (
					\Sigma_i offers[i].amount * sellPrice
				  - \Sigma_i offers[i].amount * offers[i].minPrice * buyPrice
				)
		Since smooth_mult division can be done with just a bitshift, this avoids any divisions!
		Hence, we precompute answers to queries of the form
			* \Sigma_{i : offers[i].minPrice < p} offers[i].amount
			* \Sigma_{i : offers[i].minPrice < p} offers[i].amount * offers[i].minPrice
		Note that offers[i].minPrice < sellPrice / buyPrice, for these offers.
		Hence, offers[i].minPrice * buyPrice < sellPrice
	*/

	fullPrice := ApplySmoothMult(sellPrice, smoothMult)
	partPrice := sellPrice

	fullPCS := ob.GetPCS(fullPrice, buyPrice)
	partPCS := ob.GetPCS(partPrice, buyPrice)

	fullAmt := fullPCS.CumForSale
	partAmt := partPCS.CumForSale - fullAmt

	fullAmtTimesPrice := fullPCS.CumForSaleTimesPrice
	partAmtTimesPrice := partPCS.CumForSaleTimesPrice - fullAmtTimesPrice

	// \Sigma_i offers[i].amount * sellPrice
	partAmtTimesSellPrice := partAmt * sellPrice
	// \Sigma_i offers[i].amount * offers[i].minPrice * buyPrice
	partAmtTimesMinPriceTimesBuyPrice := partAmtTimesPrice * buyPrice
	// \Sigma_i offers[i].amount * sellPrice - \Sigma_i offers[i].amount * offers[i].minPrice * buyPrice
	partAmtSum := partAmtTimesSellPrice - partAmtTimesMinPriceTimesBuyPrice
	// 2^{smoothMult} * (\Sigma_i offers[i].amount * sellPrice - \Sigma_i offers[i].amount * offers[i].minPrice * buyPrice)
	partSold := partAmtSum * math.Pow(2, float64(smoothMult))
	fullSold := fullAmt * sellPrice

	return fullSold + partSold
}
