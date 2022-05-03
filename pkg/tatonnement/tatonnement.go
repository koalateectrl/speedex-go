package tatonnement

import (
	"fmt"
	"math"

	"github.com/sandymule/speedex-go/pkg/assets"
	"github.com/sandymule/speedex-go/pkg/orderbook"
)

type Oracle struct {
	ObManager orderbook.Manager
}

func (to Oracle) ComputePrices(params CtrlParams, prices map[assets.Asset]float64, printFreq uint32) map[assets.Asset]float64 {
	ctrlParams := NewCtrlParamsWrapper(params)
	fmt.Println(to)
	baseSpyDmd := to.ObManager.SpyDmdQuery(prices, ctrlParams.MParams.MSmoothMult)
	baseObj := baseSpyDmd.GetObj()
	fmt.Println(baseSpyDmd)

	stepSize := ctrlParams.KInitStepSize

	for !ctrlParams.Done() {
		ctrlParams.IncrRnd()
		trialPrices := ctrlParams.SetTrialPrices(prices, baseSpyDmd, stepSize)
		trialSpyDmd := to.ObManager.SpyDmdQuery(trialPrices, ctrlParams.MParams.MSmoothMult)
		trialObj := trialSpyDmd.GetObj()

		if trialObj.Val <= baseObj.Val || stepSize < ctrlParams.KMinStepSize {
			prices = trialPrices
			baseSpyDmd = trialSpyDmd
			baseObj = trialObj
			stepSize = ctrlParams.StepUp(math.Max(stepSize, ctrlParams.KMinStepSize))
		} else {
			stepSize = ctrlParams.StepDown(stepSize)
		}

		if printFreq > 0 && ctrlParams.MRnd%printFreq == 0 {
			fmt.Println("----------------------------------------------------")
			fmt.Printf("TATONNEMENT STEP: step size: %f round number: %d\n", stepSize, ctrlParams.MRnd)
			for asset, price := range prices {
				dmd := baseSpyDmd.GetDelta(asset)
				fmt.Printf("TATONNEMENT: %s, Price: %f, Demand: %f\n", string(asset), price, dmd)
			}
		}
	}

	return prices
}
