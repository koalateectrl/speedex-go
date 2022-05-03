package dmdutils

import (
	"testing"

	"github.com/sandymule/speedex-go/pkg/assets"
	"github.com/stretchr/testify/assert"
)

func TestAddSpyDmd(t *testing.T) {
	actualOutput := make(SpyDmd)
	actualOutput.AddSpyDmd("ETH", "USDT", 50)

	expectedOutput := make(SpyDmd)
	expectedOutput["ETH"] = SpyDmdPair{50, 0}
	expectedOutput["USDT"] = SpyDmdPair{0, 50}

	assert.Equal(t, expectedOutput, actualOutput, "Maps (when adding new key) are not equal.")

	actualOutput.AddSpyDmd("ETH", "USDT", 30)
	expectedOutput["ETH"] = SpyDmdPair{80, 0}
	expectedOutput["USDT"] = SpyDmdPair{0, 80}

	assert.Equal(t, expectedOutput, actualOutput, "Maps (when adding to existing key) are not equal.")

	actualOutput.AddSpyDmd("USDT", "ETH", 45)
	expectedOutput["ETH"] = SpyDmdPair{80, 45}
	expectedOutput["USDT"] = SpyDmdPair{45, 80}

	assert.Equal(t, expectedOutput, actualOutput, "Maps (rev when adding to existing key) are not equal.")
}

func TestAddSpyDmdPair(t *testing.T) {
	actualOutput := make(SpyDmd)
	actualOutput.AddSpyDmdPair(assets.AssetPair{Buy: "USDT", Sell: "ETH"}, 50)

	expectedOutput := make(SpyDmd)
	expectedOutput["ETH"] = SpyDmdPair{50, 0}
	expectedOutput["USDT"] = SpyDmdPair{0, 50}

	assert.Equal(t, expectedOutput, actualOutput, "Maps (when adding new key) are not equal.")

	actualOutput.AddSpyDmdPair(assets.AssetPair{Buy: "USDT", Sell: "ETH"}, 30)
	expectedOutput["ETH"] = SpyDmdPair{80, 0}
	expectedOutput["USDT"] = SpyDmdPair{0, 80}

	assert.Equal(t, expectedOutput, actualOutput, "Maps (when adding to existing key) are not equal.")

	actualOutput.AddSpyDmdPair(assets.AssetPair{Buy: "ETH", Sell: "USDT"}, 45)
	expectedOutput["ETH"] = SpyDmdPair{80, 45}
	expectedOutput["USDT"] = SpyDmdPair{45, 80}

	assert.Equal(t, expectedOutput, actualOutput, "Maps (rev when adding to existing key) are not equal.")

}

func TestGetDelta(t *testing.T) {
	actualInput := make(SpyDmd)
	actualInput.AddSpyDmd("ETH", "USDT", 50)
	actualInput.AddSpyDmd("USDT", "ETH", 90)

	actualOutput1 := actualInput.GetDelta("ETH")
	expectedOutput1 := 40.0

	assert.Equal(t, expectedOutput1, actualOutput1, "ETH Deltas are not equal.")

	actualOutput2 := actualInput.GetDelta("USDT")
	expectedOutput2 := -40.0

	assert.Equal(t, expectedOutput2, actualOutput2, "ETH Deltas are not equal.")
}

func TestGetObj(t *testing.T) {
	actualInput := make(SpyDmd)
	actualInput.AddSpyDmdPair(assets.AssetPair{Buy: "USDT", Sell: "ETH"}, 2400)
	actualInput.AddSpyDmdPair(assets.AssetPair{Buy: "ETH", Sell: "USDT"}, 2000)
	actualOutput := actualInput.GetObj()

	expectedOutput := 320000

	assert.EqualValues(t, expectedOutput, actualOutput.Val, "Objective Function values are not equal.")

}
