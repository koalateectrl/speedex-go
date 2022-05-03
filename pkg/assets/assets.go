package assets

import "fmt"

type Asset string

type AssetPair struct {
	Buy  Asset
	Sell Asset
}

/*
type Price struct {
	N float64 // numerator
	D float64 // denominator
} */

func (ap AssetPair) String() string {
	return fmt.Sprintf("(%s / %s)", string(ap.Buy), string(ap.Sell))
}
