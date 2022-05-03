package tatonnement

import (
	"testing"

	"github.com/sandymule/speedex-go/pkg/assets"
	"github.com/sandymule/speedex-go/pkg/orderbook"
	"github.com/stretchr/testify/assert"
)

func TestIncrRound(t *testing.T) {
	ctrlParams := NewCtrlParamsWrapper(
		CtrlParams{MSmoothMult: 5, MMaxRnds: 50,
			MStepUp: 40, MStepDown: 25, MStepSizeRdx: 5, MStepRdx: 30})

	prevRnd := ctrlParams.MRnd
	ctrlParams.IncrRnd()
	newRnd := ctrlParams.MRnd

	assert.Equal(t, prevRnd+1, newRnd, "Increment Round Failed.")
}

func TestDone(t *testing.T) {
	ctrlParams := NewCtrlParamsWrapper(
		CtrlParams{MSmoothMult: 5, MMaxRnds: 50,
			MStepUp: 40, MStepDown: 25, MStepSizeRdx: 5, MStepRdx: 30})

	ctrlParams.MRnd = 49
	actualOutput := ctrlParams.Done()
	expectedOutput := false

	assert.Equal(t, expectedOutput, actualOutput, "Done should be false.")

	ctrlParams.MRnd = 50
	actualOutput = ctrlParams.Done()
	expectedOutput = true

	assert.Equal(t, expectedOutput, actualOutput, "Done should be true.")
}

func TestImposePriceBounds(t *testing.T) {
	ctrlParams := NewCtrlParamsWrapper(
		CtrlParams{MSmoothMult: 5, MMaxRnds: 50,
			MStepUp: 40, MStepDown: 25, MStepSizeRdx: 5, MStepRdx: 30})

	actualOutput := ctrlParams.ImposePriceBounds(0)
	var expectedOutput float64 = 1.0

	assert.EqualValues(t, expectedOutput, actualOutput, "Below minimum price.")

	actualOutput = ctrlParams.ImposePriceBounds(100000000000000000000000000000000000000000)
	expectedOutput = ctrlParams.KPriceMax

	assert.EqualValues(t, expectedOutput, actualOutput, "Above maximum price.")

	actualOutput = ctrlParams.ImposePriceBounds(796)
	expectedOutput = 796.0

	assert.EqualValues(t, expectedOutput, actualOutput, "Valid candidate price.")

}

func TestStepUp(t *testing.T) {
	ctrlParams := NewCtrlParamsWrapper(
		CtrlParams{MSmoothMult: 5, MMaxRnds: 50,
			MStepUp: 40, MStepDown: 25, MStepSizeRdx: 5, MStepRdx: 30})

	actualOutput := ctrlParams.StepUp(2)
	expectedOutput := 2.5

	assert.EqualValues(t, expectedOutput, actualOutput, "Step Up does not match.")

}

func TestStepDown(t *testing.T) {
	ctrlParams := NewCtrlParamsWrapper(
		CtrlParams{MSmoothMult: 5, MMaxRnds: 50,
			MStepUp: 40, MStepDown: 25, MStepSizeRdx: 5, MStepRdx: 30})

	actualOutput := ctrlParams.StepDown(8)
	expectedOutput := 6.25

	assert.EqualValues(t, expectedOutput, actualOutput, "Step Down does not match.")

}

func TestSetTrialPrice(t *testing.T) {
	ctrlParams := NewCtrlParamsWrapper(
		CtrlParams{MSmoothMult: 5, MMaxRnds: 50,
			MStepUp: 40, MStepDown: 25, MStepSizeRdx: 5, MStepRdx: 30})

	curPrice := 4500.0
	excessDmd := -400.0
	stepSize := ctrlParams.KInitStepSize

	actualOutput := ctrlParams.SetTrialPrice(curPrice, excessDmd, stepSize)
	expectedOutput := 4499.89

	assert.InEpsilon(t, expectedOutput, actualOutput, 0.1, "Set Trial with negative excess demand does not match.")

	curPrice = 4500.0
	excessDmd = 400.0
	stepSize = ctrlParams.KInitStepSize

	actualOutput = ctrlParams.SetTrialPrice(curPrice, excessDmd, stepSize)
	expectedOutput = 5000.11

	assert.InEpsilon(t, expectedOutput, actualOutput, 0.1, "Set Trial with positive excess demand does not match.")

}

func TestComputePrices(t *testing.T) {
	var oracle Oracle

	cp := CtrlParams{MSmoothMult: 2, MMaxRnds: 50,
		MStepUp: 40, MStepDown: 25, MStepSizeRdx: 5, MStepRdx: 30}
	prices := make(map[assets.Asset]float64)

	prices["ETH"] = 4500
	prices["USDT"] = 1

	mgr := make(orderbook.Manager)

	actualOrderbookEth := make(orderbook.Orderbook, 0)
	actualOrderbookEth = append(actualOrderbookEth, orderbook.PriceCompStats{SellPrice: 0, CumForSale: 0, CumForSaleTimesPrice: 0})
	actualOrderbookEth = append(actualOrderbookEth, orderbook.PriceCompStats{SellPrice: 0.0002, CumForSale: 5000, CumForSaleTimesPrice: 1})
	actualOrderbookEth = append(actualOrderbookEth, orderbook.PriceCompStats{SellPrice: 0.00025, CumForSale: 6000, CumForSaleTimesPrice: 1.5})
	mgr[assets.AssetPair{Buy: assets.Asset("ETH"), Sell: assets.Asset("USDT")}] = actualOrderbookEth

	actualOrderbookUsdt := make(orderbook.Orderbook, 0)
	actualOrderbookUsdt = append(actualOrderbookUsdt, orderbook.PriceCompStats{SellPrice: 0, CumForSale: 0, CumForSaleTimesPrice: 0})
	actualOrderbookUsdt = append(actualOrderbookUsdt, orderbook.PriceCompStats{SellPrice: 4200, CumForSale: 2, CumForSaleTimesPrice: 8400})
	mgr[assets.AssetPair{Buy: assets.Asset("USDT"), Sell: assets.Asset("ETH")}] = actualOrderbookUsdt

	oracle.ObManager = mgr

	actualPrices := oracle.ComputePrices(cp, prices, 1)
	actualOutput := oracle.ObManager.SpyDmdQuery(actualPrices, cp.MSmoothMult)

	expectedOutput := 2141.278
	assert.InDelta(t, expectedOutput, actualOutput["ETH"].Spy, 1, "Supply does not match expected.")

	for asset := range actualOutput {
		assert.InDelta(t, actualOutput[asset].Spy, actualOutput[asset].Dmd, 1, "Supply does not match demand")
	}
}
