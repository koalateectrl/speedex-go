package dmdutils

import (
	"math"

	"github.com/sandymule/speedex-go/pkg/assets"
)

type SpyDmdPair struct {
	Spy float64
	Dmd float64
}

type SpyDmd map[assets.Asset]SpyDmdPair

type ObjFunc struct {
	Val float64
}

func (sdp SpyDmdPair) IncrAmt(amt float64, isSpy bool) SpyDmdPair {
	if isSpy {
		sdp.Spy += amt
	} else {
		sdp.Dmd += amt
	}
	return sdp
}

func (sd SpyDmd) AddSpyDmdPair(tradingPair assets.AssetPair, amt float64) {
	sd.AddSpyDmd(tradingPair.Sell, tradingPair.Buy, amt)
}

func (sd SpyDmd) AddSpyDmd(sell assets.Asset, buy assets.Asset, amt float64) {

	if _, ok := sd[sell]; ok {
		sd[sell] = sd[sell].IncrAmt(amt, true)
	} else {
		sd[sell] = SpyDmdPair{amt, 0}
	}

	if _, ok := sd[buy]; ok {
		sd[buy] = sd[buy].IncrAmt(amt, false)
	} else {
		sd[buy] = SpyDmdPair{0, amt}
	}
}

func (sd SpyDmd) GetDelta(asset assets.Asset) float64 {
	return sd[asset].Dmd - sd[asset].Spy
}

func (sd SpyDmd) GetObj() *ObjFunc {
	return NewObjFunc(sd)
}

func NewObjFunc(sd SpyDmd) *ObjFunc {
	tof := new(ObjFunc)
	for _, onesd := range sd {
		tof.Val += math.Pow(onesd.Dmd-onesd.Spy, 2)
	}

	return tof
}
