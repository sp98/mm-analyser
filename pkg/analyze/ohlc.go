package analyze

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const (
	bullshMarubuzoAfterDecline = "BullishMarubuzoAfterDecline"
	openLowHigh                = "OpenLowHigh"
	doziAfterDecline           = "DoziAfterDecline"
	bullishHammerAfterDecline  = "BullishHammerAfterDecline"
	bearishHammerAfterDecline  = "BearishHammerAfterDecline"
	endOfDecline               = "EndOfDecline"
	bearishMarubuzoAfterRally  = "BearishMarubuzoAfterRally"
	doziAfterRally             = "DoziAfterRally"
	shootingStarAfterDecline   = "ShootingStarAfterDecline"
	shootingStartAfterRally    = "ShootingStartAfterRally"
	endOfRally                 = "EndOfRally"
)

func (i *Instrument) ohlcAnalyser(result *Result) {

	ohlc := *i.OHLC
	shortTrend, _ := getShortTermTrend(ohlc[1:])
	if shortTrend == "" {
		log.Printf("No short term trend observed in the Instrument %s", i.Name)
		//return
	}

	//Uptrend Indicators
	isBullishMaru := isBullishMarubuzo(ohlc[0])
	if shortTrend == "decline" && isBullishMaru {
		result.UpdateResult(bullshMarubuzoAfterDecline, i)
		return
	}

	isDozi := isDozi(ohlc[0])
	if shortTrend == "decline" && isDozi {
		result.UpdateResult(doziAfterDecline, i)
		return
	}

	isbullishHammer := isBullishHammer(ohlc[0])
	if shortTrend == "decline" && isbullishHammer {
		result.UpdateResult(bullishHammerAfterDecline, i)
		return
	}

	isbearishHammer := isBearishHammer(ohlc[0])
	if shortTrend == "decline" && isbearishHammer {
		result.UpdateResult(bearishHammerAfterDecline, i)
		return
	}

	//Downtrend Indicators
	isBearishMaru := isBearishMarubuzo(ohlc[0])
	if shortTrend == "rally" && isBearishMaru {
		result.UpdateResult(bearishMarubuzoAfterRally, i)
		return
	}

	if shortTrend == "rally" && isDozi {
		result.UpdateResult(bullshMarubuzoAfterDecline, i)
		return
	}

	isinvertedHammer := isInvertedHammer(ohlc[0])
	if shortTrend == "rally" && isinvertedHammer {
		result.UpdateResult(shootingStartAfterRally, i)
		return
	}

	if shortTrend == "rally" && isDozi {
		result.UpdateResult(doziAfterRally, i)
		return
	}

}

//UpdateResult updates the analysis result
func (r *Result) UpdateResult(resultType string, i *Instrument) {
	r.Mux.Lock()
	defer r.Mux.Unlock()

	switch resultType {
	case openLowHigh:
		r.OpenLowHigh = append(r.OpenLowHigh, *i)
	case bullshMarubuzoAfterDecline:
		r.BullishMarubuzoAfterDecline = append(r.BullishMarubuzoAfterDecline, *i)
		break
	case doziAfterDecline:
		r.DoziAfterDecline = append(r.DoziAfterDecline, *i)
		break
	case bullishHammerAfterDecline:
		r.BullishHammerAfterDecline = append(r.BullishHammerAfterDecline, *i)
		break
	case bearishHammerAfterDecline:
		r.BearishHammerAfterDecline = append(r.BearishHammerAfterDecline, *i)
		break
	case endOfDecline:
		r.EndOfDecline = append(r.EndOfDecline, *i)
		break
	case bearishMarubuzoAfterRally:
		r.BearishMarubuzoAfterRally = append(r.BearishMarubuzoAfterRally, *i)
		break
	case doziAfterRally:
		r.DoziAfterRally = append(r.DoziAfterRally, *i)
		break
	case shootingStarAfterDecline:
		r.ShootingStarAfterDecline = append(r.ShootingStarAfterDecline, *i)
		break
	case shootingStartAfterRally:
		r.ShootingStartAfterRally = append(r.ShootingStartAfterRally, *i)
		break
	case endOfRally:
		r.EndOfRally = append(r.EndOfRally, *i)
		break
	}

}

//GetOHLC fetches the OHLC data for an instruement form the TickStore Rest API
func GetOHLC(token, interval string) (*[]OHLC, error) {
	url := fmt.Sprintf(DailyOHCLAPI, os.Getenv("TICK_STORE_API"), token, interval)
	resp, err := getWithAuth(url, os.Getenv("API_USER_NAME"), os.Getenv("API_PASSWORD"))
	if err != nil {
		return nil, fmt.Errorf("error fetching ohlc. %+v", err)
	}
	var result []OHLC
	json.NewDecoder(resp.Body).Decode(&result)
	return &result, nil

}

//StoreOHLCResult calls the API end point to store the OHLC analysis results
func StoreOHLCResult(interval string, res *Result) error {
	url := fmt.Sprintf(ResultStoreAPI, os.Getenv("RESULT_STORE_API"), interval)
	_, err := postWithAuth(url, os.Getenv("API_USER_NAME"), os.Getenv("API_PASSWORD"), res)
	if err != nil {
		return fmt.Errorf("error storing ohlc analysis result. %+v", err)
	}
	return nil
}

//GetLastesStockData get latest stock data
func GetLastesStockData(token string) (map[string]StockData, error) {
	url := fmt.Sprintf(LatestStockDataAPI, os.Getenv("TICK_STORE_API"), token)
	log.Println("SP - url - ", url)
	resp, err := getWithAuth(url, os.Getenv("API_USER_NAME"), os.Getenv("API_PASSWORD"))
	if err != nil {
		return nil, fmt.Errorf("error fetching ohlc. %+v", err)
	}

	log.Printf("Body %+v", resp.Body)
	var result map[string]StockData
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil

}

func getOpenLowHigh(token string) []Instrument {

	olhInsturmentList := []Instrument{}
	resMap, _ := GetLastesStockData(token)

	for _, res := range resMap {
		if res.Open == res.Low || res.Open == res.High {
			insturment := NewInsturment(res.Name, res.Symnbol, res.Token, res.Exchange, nil)
			olhInsturmentList = append(olhInsturmentList, *insturment)
		}
	}

	return olhInsturmentList

}

//Gives the trend before the current Candlestick pattern
func getShortTermTrend(ohlcList []OHLC) (string, int) {
	trend := ""
	trendCount := 0

	if len(ohlcList) < 3 {
		return trend, trendCount
	}

	//Last three candles should make higher highs and lower lows for Rally
	for i := 0; i < len(ohlcList)-1; i++ {
		if ohlcList[i].High > ohlcList[i+1].High && ohlcList[i].Low > ohlcList[i+1].Low { //Todo: should equality also be used.
			trendCount = trendCount + 1
			continue
		}
		//If updtrend count is >=3, then consider it as a rally
		if trendCount >= 3 {
			trend = "rally"
			return trend, trendCount
		}

	}

	//Reintialize trendcount back to 0
	trendCount = 0

	//Last three candles should make lower highs and higher lows for decline
	for i := 0; i < len(ohlcList)-1; i++ {
		if ohlcList[i].High < ohlcList[i+1].High && ohlcList[i].Low < ohlcList[i+1].Low {
			trendCount = trendCount + 1
			continue
		}
		//If downtrend count is >=3, then consider it as a decline
		if trendCount >= 3 {
			trend = "decline"
			return trend, trendCount
		}
	}

	//No trend found and trend count should be 0
	return trend, 0

}

func isBullish(ohlcList []OHLC) (bool, int) {

	isBull := true
	lastCandleStick := ohlcList[0]
	if lastCandleStick.Open >= lastCandleStick.Close {
		isBull = false
		return isBull, 0
	}

	var trendCount = 1
	for _, ohlc := range ohlcList {
		if ohlc.Open > ohlc.Close {
			break
		}
		trendCount = trendCount + 1
	}

	return isBull, trendCount

}

func isBearish(csDetails []OHLC) (bool, int) {

	isBear := true
	lastCandleStick := csDetails[0]
	if lastCandleStick.Open < lastCandleStick.Close {
		isBear = false
		return isBear, 0
	}
	var trendCount = 1
	for _, cs := range csDetails {
		if cs.Open <= cs.Close {
			break
		}
		trendCount = trendCount + 1
	}

	return isBear, trendCount

}

func isBullishMarubuzo(csDetails OHLC) bool {
	if csDetails.Open < csDetails.Close {
		if csDetails.Open == csDetails.Low && csDetails.Close == csDetails.High {
			return true
		}

		//Candlestick with body > 80% of the today candlestick size
		if (((csDetails.Close - csDetails.Open) / (csDetails.High - csDetails.Low)) * 100) > 80 {
			return true
		}
	}

	return false

}

func isBearishMarubuzo(ohlc OHLC) bool {

	if ohlc.Open > ohlc.Close {
		if ohlc.Open == ohlc.High && ohlc.Close == ohlc.Low {
			return true
		}
		//Candlestick with body > 80% of the today candlestick size
		if (((ohlc.Open - ohlc.Close) / (ohlc.High - ohlc.Low)) * 100) > 80 {
			return true
		}
	}

	return false

}

func isDozi(ohlc OHLC) bool {
	if ohlc.Open == ohlc.Close && (ohlc.High != ohlc.Open || ohlc.Low != ohlc.Open) {
		return true
	}

	return false

}

func isInvertedHammer(ohlc OHLC) bool {
	if ohlc.Open < ohlc.Close {
		if (2*(ohlc.Close-ohlc.Open) < (ohlc.High - ohlc.Close)) && ((ohlc.Open - ohlc.Low) < (ohlc.Close - ohlc.Open)) {
			return true
		}
	} else if ohlc.Open > ohlc.Close {
		if (2*(ohlc.Open-ohlc.Close) < (ohlc.High - ohlc.Open)) && ((ohlc.Close - ohlc.Low) < (ohlc.Open - ohlc.Close)) {
			return true
		}
	}

	return false
}

func isBullishHammer(ohlc OHLC) bool {
	if ohlc.Open < ohlc.Close {
		//Shadow twice the length body and high == close or very small shadow on the top
		if ((ohlc.Open - ohlc.Low) >= 2*(ohlc.Close-ohlc.Open)) &&
			((ohlc.High == ohlc.Close) || ((ohlc.High - ohlc.Close) < (ohlc.Close - ohlc.Open))) {
			return true
		}
	}

	return false
}

func isBearishHammer(ohlc OHLC) bool {
	if ohlc.Open > ohlc.Close {
		//Shadow twice the length body and high == open or very small shadow on the top
		if ((ohlc.Close - ohlc.Low) >= 2*(ohlc.Open-ohlc.Close)) &&
			((ohlc.High == ohlc.Open) || ((ohlc.High - ohlc.Open) < (ohlc.Open - ohlc.Close))) {
			return true
		}
	}
	return false
}

//isLowerHighsEngulfingPatter checks for patter where lower highs are made but lows may be lower or higher (making the previous pattern engulfuing the previous one)
func lowerHighsEngulfingPatternCount(ohlc []OHLC) int {
	count := 0
	for i := 0; i < len(ohlc)-1; i++ {
		if ohlc[i].High < ohlc[i+1].High && ((ohlc[i].Low < ohlc[i+1].Low) || (ohlc[i].Low > ohlc[i+1].Low)) {
			count = count + 1
			continue
		}
		return count
	}
	return count
}

func higherLowsEngulfingPatternCount(ohlc []OHLC) int {
	count := 0
	for i := 0; i < len(ohlc)-1; i++ {
		if ohlc[i].Low > ohlc[i+1].Low && ((ohlc[i].High > ohlc[i+1].High) || (ohlc[i].High <= ohlc[i+1].High)) {
			count = count + 1
			continue
		}
		return count
	}
	return count
}
