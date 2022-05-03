package orderbook

import (
	"testing"

	"github.com/sandymule/speedex-go/pkg/assets"
	"github.com/sandymule/speedex-go/pkg/dmdutils"
	"github.com/stretchr/testify/assert"
)

func TestApplySmoothMult(t *testing.T) {
	actualInput := 5.0
	smoothMult := 0
	actualOutput := ApplySmoothMult(actualInput, uint8(smoothMult))
	expectedOutput := 5.0

	assert.Equal(t, expectedOutput, actualOutput, "No Smooth Mult")

	smoothMult = 5
	actualOutput = ApplySmoothMult(actualInput, uint8(smoothMult))
	expectedOutput = 4.84375

	assert.Equal(t, expectedOutput, actualOutput, "With Smooth Mult")
}

func TestGetPCS(t *testing.T) {
	actualInput := make(Orderbook, 0)
	actualInput = append(actualInput, PriceCompStats{SellPrice: 0, CumForSale: 0, CumForSaleTimesPrice: 0})
	actualOutput1 := actualInput.GetPCS(5, 4)
	expectedOutput1 := PriceCompStats{SellPrice: 0, CumForSale: 0, CumForSaleTimesPrice: 0}

	assert.Equal(t, expectedOutput1, actualOutput1, "Slices (only zero element) are not equal.")

	actualInput = append(actualInput, PriceCompStats{SellPrice: 0.0002, CumForSale: 39980, CumForSaleTimesPrice: 7.996})
	actualInput = append(actualInput, PriceCompStats{SellPrice: 0.0006, CumForSale: 47982, CumForSaleTimesPrice: 12.7972})
	actualOutput2 := actualInput.GetPCS(1, 10000)
	expectedOutput2 := PriceCompStats{SellPrice: 0, CumForSale: 0, CumForSaleTimesPrice: 0}

	assert.Equal(t, expectedOutput2, actualOutput2, "Slices (below any offer price) are not equal.")

	actualOutput3 := actualInput.GetPCS(3, 10000)
	expectedOutput3 := PriceCompStats{SellPrice: 0.0002, CumForSale: 39980, CumForSaleTimesPrice: 7.996}

	assert.Equal(t, expectedOutput3, actualOutput3, "Slices (middle of offer prices) are not equal.")

	actualOutput4 := actualInput.GetPCS(7, 10000)
	expectedOutput4 := PriceCompStats{SellPrice: 0.0006, CumForSale: 47982, CumForSaleTimesPrice: 12.7972}

	assert.Equal(t, expectedOutput4, actualOutput4, "Slices (above all offer prices) are not equal.")
}

func TestGetAmt(t *testing.T) {
	ob := make(Orderbook, 0)
	ob = append(ob, PriceCompStats{SellPrice: 0, CumForSale: 0, CumForSaleTimesPrice: 0})
	ob = append(ob, PriceCompStats{SellPrice: 4200, CumForSale: 2, CumForSaleTimesPrice: 8400})
	smoothMult := 2

	actualOutput := ob.GetAmt(8000, 1, uint8(smoothMult))
	expectedOutput := 16000

	assert.EqualValues(t, expectedOutput, actualOutput, "Only full sales.")

	actualOutput = ob.GetAmt(4500, 1, uint8(smoothMult))
	expectedOutput = 2400

	assert.EqualValues(t, expectedOutput, actualOutput, "Only partial sales.")

	actualOutput = ob.GetAmt(0, 1, uint8(smoothMult))
	expectedOutput = 0

	assert.EqualValues(t, expectedOutput, actualOutput, "No sales.")

	ob = append(ob, PriceCompStats{SellPrice: 7000, CumForSale: 5, CumForSaleTimesPrice: 29400})
	actualOutput = ob.GetAmt(8000, 1, uint8(smoothMult))
	expectedOutput = 28000

	assert.EqualValues(t, expectedOutput, actualOutput, "Full and partial sales.")

}

func TestDmdQuery(t *testing.T) {
	actualInput := make(Manager)
	actualPrices := make(map[assets.Asset]float64)
	actualPrices["ETH"] = 4500
	actualPrices["USDT"] = 1

	actualOrderbookEth := make(Orderbook, 0)
	actualOrderbookEth = append(actualOrderbookEth, PriceCompStats{SellPrice: 0, CumForSale: 0, CumForSaleTimesPrice: 0})
	actualOrderbookEth = append(actualOrderbookEth, PriceCompStats{SellPrice: 0.0002, CumForSale: 5000, CumForSaleTimesPrice: 1})
	actualOrderbookEth = append(actualOrderbookEth, PriceCompStats{SellPrice: 0.00025, CumForSale: 6000, CumForSaleTimesPrice: 1.5})
	actualInput[assets.AssetPair{Buy: assets.Asset("ETH"), Sell: assets.Asset("USDT")}] = actualOrderbookEth

	actualOrderbookUsdt := make(Orderbook, 0)
	actualOrderbookUsdt = append(actualOrderbookUsdt, PriceCompStats{SellPrice: 0, CumForSale: 0, CumForSaleTimesPrice: 0})
	actualOrderbookUsdt = append(actualOrderbookUsdt, PriceCompStats{SellPrice: 4200, CumForSale: 2, CumForSaleTimesPrice: 8400})
	actualInput[assets.AssetPair{Buy: assets.Asset("USDT"), Sell: assets.Asset("ETH")}] = actualOrderbookUsdt

	actualOutput := actualInput.SpyDmdQuery(actualPrices, 2)

	expectedOutput := make(dmdutils.SpyDmd)
	expectedOutput["ETH"] = dmdutils.SpyDmdPair{Spy: 2400, Dmd: 2000}
	expectedOutput["USDT"] = dmdutils.SpyDmdPair{Spy: 2000, Dmd: 2400}

	assert.Equal(t, expectedOutput, actualOutput, "Maps of Supply Demand are not equal.")
}
