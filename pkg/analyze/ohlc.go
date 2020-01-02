package analyze

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var (
	//FromTime queries influx db from specific time period.
	FromTime = ""
	//ToTime queries influx db till specific time period.
	ToTime = ""
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
	bullishMarubuzo            = "BullishMarubuzo"
	bearishMarubuzo            = "BearishMarubuzo"
	dozi                       = "Dozi"
	hammer                     = "Hammer"
	shootingStar               = "ShootingStar"
)

func (i *Instrument) ohlcAnalyser(result *Result) {

	ohlc := *i.OHLC
	shortTrend, count := getShortTermTrend(ohlc[1:])
	if shortTrend == "" {
		log.Printf("no short term ternd observed in the instrument %q", i.Name)
	} else {
		log.Printf("short term trend of %q with count %d is observed in the instrument %q", shortTrend, count, i.Name)
	}

	//Uptrend Indicators
	isBullishMaru := isBullishMarubuzo(ohlc[0])
	if shortTrend == "decline" && count >= 3 && isBullishMaru {
		result.UpdateResult(bullshMarubuzoAfterDecline, i)

	}

	isDozi := isDozi(ohlc[0])
	if shortTrend == "decline" && count >= 3 && isDozi {
		result.UpdateResult(doziAfterDecline, i)

	}

	isbullishHammer := isBullishHammer(ohlc[0])
	if shortTrend == "decline" && count >= 3 && isbullishHammer {
		result.UpdateResult(bullishHammerAfterDecline, i)

	}

	isbearishHammer := isBearishHammer(ohlc[0])
	if shortTrend == "decline" && count >= 3 && isbearishHammer {
		result.UpdateResult(bearishHammerAfterDecline, i)

	}

	if hasDeclineEnded(shortTrend, count, ohlc[0:2]) {
		result.UpdateResult(endOfDecline, i)

	}

	//Downtrend Indicators
	isBearishMaru := isBearishMarubuzo(ohlc[0])
	if shortTrend == "rally" && count >= 3 && isBearishMaru {
		result.UpdateResult(bearishMarubuzoAfterRally, i)

	}

	isinvertedHammer := isInvertedHammer(ohlc[0])
	if shortTrend == "rally" && count >= 3 && isinvertedHammer {
		result.UpdateResult(shootingStartAfterRally, i)

	}

	if shortTrend == "rally" && count >= 3 && isDozi {
		result.UpdateResult(doziAfterRally, i)

	}

	if hasRallyEnded(shortTrend, count, ohlc[0:2]) {
		result.UpdateResult(endOfRally, i)

	}

	//Other candlestick types:
	if isDozi {
		result.UpdateResult(dozi, i)
	}

	if isBullishMaru {
		result.UpdateResult(bullishMarubuzo, i)
	}

	if isBearishMaru {
		result.UpdateResult(bearishMarubuzo, i)
	}

	if isbullishHammer {
		result.UpdateResult(hammer, i)
	}

	if isbearishHammer {
		result.UpdateResult(hammer, i)
	}

	if isinvertedHammer {
		result.UpdateResult(shootingStar, i)
	}

}

//UpdateResult updates the analysis result
func (r *Result) UpdateResult(resultType string, i *Instrument) {
	i.OHLC = nil // Don't store ohlc data for now.
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
	case bearishMarubuzo:
		r.BearishMarubuzo = append(r.BearishMarubuzo, *i)
		break
	case bullishMarubuzo:
		r.BullishMarubuzo = append(r.BullishMarubuzo, *i)
		break
	case dozi:
		r.Dozi = append(r.Dozi, *i)
		break
	case hammer:
		r.Hammer = append(r.Hammer, *i)
		break
	case shootingStar:
		r.ShootingStar = append(r.ShootingStar, *i)
		break
	}

}

//GetOHLC fetches the OHLC data for an instruement form the TickStore Rest API
func GetOHLC(token, interval string) (*[]OHLC, error) {
	var url string
	if FromTime != "" && ToTime != "" {
		url = fmt.Sprintf(SpecificOHCLAPI, os.Getenv("TICK_STORE_API"), token, interval, FromTime, ToTime)
	} else {
		url = fmt.Sprintf(DailyOHCLAPI, os.Getenv("TICK_STORE_API"), token, interval)
	}

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
	resp, err := getWithAuth(url, os.Getenv("API_USER_NAME"), os.Getenv("API_PASSWORD"))
	if err != nil {
		return nil, fmt.Errorf("error fetching ohlc. %+v", err)
	}

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

	if ohlcList[0].High > ohlcList[1].High && ohlcList[0].Low > ohlcList[1].Low {
		trend = "rally"
		for i := 0; i < len(ohlcList)-1; i++ {
			if ohlcList[i].High > ohlcList[i+1].High && ohlcList[i].Low > ohlcList[i+1].Low {
				trendCount = trendCount + 1
				continue
			}
			return trend, trendCount
		}
		return trend, trendCount

	} else if ohlcList[0].High < ohlcList[1].High && ohlcList[0].Low < ohlcList[1].Low {
		trend = "decline"
		for i := 0; i < len(ohlcList)-1; i++ {
			if ohlcList[i].High < ohlcList[i+1].High && ohlcList[i].Low < ohlcList[i+1].Low {
				trendCount = trendCount + 1
				continue
			}
			return trend, trendCount
		}
		return trend, trendCount

	}

	return "", 0
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

func hasRallyEnded(trend string, trendCount int, ohlc []OHLC) bool {
	if trend == "rally" && trendCount >= 3 {
		if len(ohlc) == 2 {
			// Rally ends if the latest candlestick is red or if  Higher high and higher lows are not made.
			if ohlc[0].Open > ohlc[0].Close || ohlc[0].High < ohlc[1].High || ohlc[0].Low < ohlc[1].High {
				return true
			}
		}
	}
	return false
}

func hasDeclineEnded(trend string, trendCount int, ohlc []OHLC) bool {
	if trend == "decline" && trendCount >= 3 {
		if len(ohlc) == 2 {
			// Decline ends if Lower high and Lower lows are not made. Or if the latest candlestick is green
			if ohlc[0].Open < ohlc[0].Close || ohlc[0].High > ohlc[1].High || ohlc[0].Low > ohlc[1].High {
				return true
			}
		}
	}
	return false
}
