package analyze

import (
	"log"
	"sync"

	"github.com/sp98/analyzer/pkg/utility"
)

var (
	//DailyOHCLAPI is the API url for fetching daily OHLC data for an instrument
	DailyOHCLAPI = "%s/v1/api/ohlc/%s/%s"

	//SpecificOHCLAPI is the API url for fetching OHLC data for an instrument from a specific time period
	SpecificOHCLAPI = "%s/v1/api/ohlc/%s/%s?from=%s&to=%s"
	//LatestStockDataAPI is the API end point to get the latest stock data
	LatestStockDataAPI = "%s/v1/api/stocks/%s"

	//ResultStoreAPI is the API end point to store the OHLC analysis results.
	ResultStoreAPI = "%s/v1/api/ohlcresult/%s"
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

	//Opending Trends
	OpenHighLow []Instrument `json:"OpenHighLow"`

	//Other chart patterns
	BullishMarubuzo []Instrument `json:"BullishMarubuzo"`
	BearishMarubuzo []Instrument `json:"BearishMarubuzo"`
	Dozi            []Instrument `json:"Dozi"`
	Hammer          []Instrument `json:"Hammer"`
	ShootingStar    []Instrument `json:"ShootingStar"`
}

//OHLC is the open, high, low and close price for an instrument.
type OHLC struct {
	Open  float64 `json:"Open"`
	High  float64 `json:"High"`
	Low   float64 `json:"Low"`
	Close float64 `json:"Close"`
}

//StockData is the struct for latest stock data
type StockData struct {
	Name              string  `json:"Name"`
	Symnbol           string  `json:"Symnbol"`
	Exchange          string  `json:"Exchange"`
	Token             string  `json:"Token"`
	LastPrice         float64 `json:"LastPrice"`
	AverageTradePrice float64 `json:"AverageTradePrice"`
	TotalBuy          float64 `json:"TotalBuy"`
	TotalSell         float64 `json:"TotalSell"`
	Open              float64 `json:"Open"`
	High              float64 `json:"High"`
	Low               float64 `json:"Low"`
	Close             float64 `json:"Close"`
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

//StartAnalysis starts the analysis of the stock
func StartAnalysis(stocks, interval string) {
	var wg sync.WaitGroup
	stockList := utility.GetStocks(stocks)
	result := &Result{Mux: &sync.Mutex{}}

	//Open Low High Analysis
	isWithInActualMarketOpenTime, _ := utility.IsWithInActualMarketOpenTime()

	if isWithInActualMarketOpenTime {
		log.Println("within actual market open time. so analyzing open low high strategy")
		OpenHighLowAnalysis(stocks, result)
	}

	for _, stock := range stockList {
		wg.Add(1)
		ohlc, err := GetOHLC(stock[2], interval)
		if err != nil {
			log.Printf("error in setup for stock - %q. %+v", stock[0], err)
		}
		insturment := NewInsturment(stock[0], stock[1], stock[2], stock[3], ohlc)

		go insturment.Analyze(result, &wg)
	}

	wg.Wait() //wait for all the go routines to finish

	log.Printf("Result - %+v", result)
	StoreOHLCResult(interval, result)
	utility.IsMarketOpen()

}

//Analyze the instrument's tick data
func (i *Instrument) Analyze(result *Result, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Analyzing the instrument - %+v ", i)
	if len(*i.OHLC) <= 3 {
		log.Printf("skip analysis for stock %q as number of candle sticks are less than 3", i.Name)
		return
	}

	i.ohlcAnalyser(result)

}

//OpenHighLowAnalysis does openhighlow analysis
func OpenHighLowAnalysis(stocks string, result *Result) {
	//Update Open Low High Results
	tokenString := utility.GetSubscriptionsString(stocks)
	olhResultList := getOpenHighLow(tokenString)
	for _, ohlcResult := range olhResultList {
		result.UpdateResult(openHighLow, &ohlcResult)
	}

}
