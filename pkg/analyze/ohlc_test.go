package analyze

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getInsturement() {

}

func TestRallyTrend(t *testing.T) {
	//result := &Result{Mux: &sync.Mutex{}}
	ohlc := &[]OHLC{}
	ohlc2 := &[]OHLC{}
	err := json.Unmarshal([]byte(RallyCount1), ohlc)
	err = json.Unmarshal([]byte(RallyCount5), ohlc2)
	if err != nil {
		t.Errorf("err enconding test data OHLC. %+v", err)
	}
	insturment := NewInsturment("UPL", "UPL", "2889473", "NSE", ohlc)
	insturment2 := NewInsturment("UPL", "UPL", "2889473", "NSE", ohlc2)

	trend, count := getShortTermTrend(*insturment.OHLC)
	assert.Equal(t, trend, "rally")
	assert.Equal(t, count, 1)

	trend, count = getShortTermTrend(*insturment2.OHLC)
	assert.Equal(t, trend, "rally")
	assert.Equal(t, count, 5)

}

func TestDeclineTrend(t *testing.T) {
	//result := &Result{Mux: &sync.Mutex{}}
	ohlc := &[]OHLC{}
	err := json.Unmarshal([]byte(DeclineTrend), ohlc)
	if err != nil {
		t.Errorf("err enconding test data OHLC. %+v", err)
	}

	insturment := NewInsturment("UPL", "UPL", "2889473", "NSE", ohlc)
	trend, count := getShortTermTrend(*insturment.OHLC)
	t.Errorf("Trend- %q. Count- %d", trend, count)

}

func TestEndOfRally(t *testing.T) {

	ohlc := &[]OHLC{}
	err := json.Unmarshal([]byte(EndOfRallyData), ohlc)
	if err != nil {
		t.Errorf("err enconding test data OHLC. %+v", err)
	}

	insturment := NewInsturment("UPL", "UPL", "2889473", "NSE", ohlc)
	data := *insturment.OHLC
	trend, count := getShortTermTrend(data[1:])
	assert.Equal(t, trend, "rally")
	assert.Equal(t, count, 4)

	endOfRally := hasRallyEnded(trend, count, data[0:2])
	assert.Equal(t, endOfRally, true)

}
