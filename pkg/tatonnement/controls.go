package tatonnement

import (
	"math"

	"github.com/sandymule/speedex-go/pkg/assets"
	"github.com/sandymule/speedex-go/pkg/dmdutils"
)

type CtrlParams struct {
	MSmoothMult                      uint8
	MMaxRnds                         uint32
	MStepUp, MStepDown, MStepSizeRdx uint8
	MStepRdx                         float64
}

type CtrlParamsWrapper struct {
	MParams       CtrlParams
	MRnd          uint32
	KPriceMin     float64
	KPriceMax     float64
	KMinStepSize  float64
	KInitStepSize float64
}

func NewCtrlParamsWrapper(cp CtrlParams) *CtrlParamsWrapper {
	cpw := new(CtrlParamsWrapper)
	cpw.MParams = cp
	cpw.KPriceMin = 1
	cpw.KPriceMax = math.MaxInt64 * math.Pow(2, -float64(cp.MSmoothMult+1))
	cpw.KMinStepSize = math.Pow(2, float64(cp.MStepSizeRdx+1))
	cpw.KInitStepSize = cpw.KMinStepSize
	return cpw
}

func (cpw *CtrlParamsWrapper) IncrRnd() {
	cpw.MRnd++
}

func (cpw CtrlParamsWrapper) Done() bool {
	return cpw.MRnd >= cpw.MParams.MMaxRnds
}

func (cpw CtrlParamsWrapper) ImposePriceBounds(candidate float64) float64 {
	if candidate > cpw.KPriceMax {
		return cpw.KPriceMax
	}
	if candidate < cpw.KPriceMin {
		return cpw.KPriceMin
	}
	return candidate
}

func (cpw CtrlParamsWrapper) StepUp(step float64) float64 {
	out := step * float64(cpw.MParams.MStepUp) * math.Pow(2, -float64(cpw.MParams.MStepSizeRdx))
	return out
}

func (cpw CtrlParamsWrapper) StepDown(step float64) float64 {
	out := step * float64(cpw.MParams.MStepDown) * math.Pow(2, -float64(cpw.MParams.MStepSizeRdx))
	return out
}

func (cpw CtrlParamsWrapper) SetTrialPrice(curPrice float64, excessDmd float64, stepSize float64) float64 {
	// set price for one asset
	stepTimesOldPrice := curPrice * stepSize

	var sign int8 = 1
	var uExcessDmd float64 = excessDmd
	if excessDmd <= 0 {
		sign = -1
		uExcessDmd = -excessDmd
	}

	product := stepTimesOldPrice * uExcessDmd

	delta := product * math.Pow(2, -float64(cpw.MParams.MStepRdx))

	var out float64
	if sign > 0 {
		out = curPrice + delta
		if out < curPrice {
			out = math.MaxFloat64 // overflow
		}
	} else {
		if curPrice > delta {
			out = curPrice - delta
		} else {
			out = 0
		}
	}
	return cpw.ImposePriceBounds(out)
}

func (cpw CtrlParamsWrapper) SetTrialPrices(curPrices map[assets.Asset]float64, dmds dmdutils.SpyDmd, stepSize float64) map[assets.Asset]float64 {
	// set prices for all assets
	pricesOut := make(map[assets.Asset]float64)
	for asset, curPrice := range curPrices {
		pricesOut[asset] = cpw.SetTrialPrice(curPrice, dmds.GetDelta(asset), stepSize)
	}

	return pricesOut
}
