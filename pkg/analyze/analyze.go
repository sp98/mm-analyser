package analyze

import (
	"log"
	"sync"
)

var (
	//DailyOHCLAPI is the API url for fetching daily OHLC data for an instrument
	DailyOHCLAPI = "%s/v1/api/ohlc/%s/%s"
)

//Instrument represents a partcular stock in BSE or NSE
type Instrument struct {
	Name     string  `json:"Name"`
	Exchange string  `json:"Exchange"`
	Symbol   string  `json:"Symbol"`
	Token    string  `json:"Token"`
	OHLC     *[]OHLC `json:"OHLC"`
}

//Result of the OHLC analysis
type Result struct {
	Mux *sync.Mutex
	//Uptrend indicators
	BullishMarubuzoAfterDecline []Instrument `json:"BullishMarubuzoAfterDecline"`
	DoziAfterDecline            []Instrument `json:"DoziAfterDecline"`
	BullishHammerAfterDecline   []Instrument `json:"BullishHammerAfterDecline"`
	BearishHammerAfterDecline   []Instrument `json:"BearishHammerAfterDecline"`
	EndOfDecline                []Instrument `json:"EndOfDecline"`

	//Downtrend Indicators
	BearishMarubuzoAfterRally []Instrument `json:"BearishMarubuzoAfterRally"`
	DoziAfterRally            []Instrument `json:"DoziAfterRally"`
	ShootingStarAfterDecline  []Instrument `json:"ShootingStarAfterDecline"`
	ShootingStartAfterRally   []Instrument `json:"ShootingStartAfterRally"`
	EndOfRally                []Instrument `json:"EndofRally"`

	//Others
	OpenLowHigh []Instrument `json:"OpenLowHigh"`
}

//OHLC is the open, high, low and close price for an instrument.
type OHLC struct {
	Open  float64 `json:"Open"`
	High  float64 `json:"High"`
	Low   float64 `json:"Low"`
	Close float64 `json:"Close"`
}

//NewInsturment creats a new instrument
func NewInsturment(name, symbol, token, exchange string, ohlc *[]OHLC) *Instrument {
	return &Instrument{
		Name:     name,
		Symbol:   symbol,
		Token:    token,
		Exchange: exchange,
		OHLC:     ohlc,
	}
}

//Analyze the instrument's tick data
func (i *Instrument) Analyze(result *Result) {
	log.Printf("Analyzing the instrument - %+v ", i)
	log.Printf("Instrument OHLC - %+v", i.OHLC)

	i.ohlcAnalyser(result)

}
