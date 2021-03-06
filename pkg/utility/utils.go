package utility

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	MarketOpenTime             = "%s 9:00:00"
	MarketCloseTime            = "%s 15:30:00"
	ActualMarketOpenTime       = "%s 09:30:00"
	PreMarketOpenAnalysisTime1 = "%s 9:13:00"
	PreMarketOpenAnalysisTime2 = "%s 9:15:00"
	TstringFormat              = "2006-01-02 15:04:05"
	LayOut                     = "2006-01-02 15:04:05"
	InfluxLayout               = "2006-01-02T15:04:05Z"
)

//IsMarketOpen returns true if market is open or viceversa
func IsMarketOpen(interval string) {

	marketOpenTime := getMarketOpenTime(interval)
	for {
		t, _ := IsWithInTimeRange(marketOpenTime, MarketCloseTime)
		if t {
			log.Println(" withhin market open time.")
			break
		}
		log.Println("not withhin market open time. Waiting...")
		time.Sleep(10 * time.Second)
	}

}

func getMarketOpenTime(i string) string {
	switch i {
	case "5m":
		return "%s 9:00:00"

	case "10m":
		return "%s 9:15:00"

	case "15m":
		return "%s 9:15:00"

	default:
		return "%s 9:00:00"
	}
}

//WaitBeforeAnalysis returns the number of seconds we need to wait before the next analysis to start
func WaitBeforeAnalysis(interval int) int {
	ct, _ := parseTime(LayOut, time.Now().Format(TstringFormat))
	_, min, sec := ct.Clock()
	next := min + (interval - min%interval)
	waitTime := next - min
	if waitTime%interval == 0 {
		return 0
	}
	waitTimeSeconds := (waitTime * 60) - sec
	return waitTimeSeconds

}

//IsWithInMarketOpenTime tells whether current time is withing market time and not on weekends
func IsWithInMarketOpenTime() (bool, error) {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	motString := fmt.Sprintf(MarketOpenTime, time.Now().Format("2006-01-02"))
	mot, err := time.ParseInLocation("2006-01-02 15:04:05", motString, loc)
	if err != nil {
		return false, fmt.Errorf("error parsing market open time. %+v", err)
	}

	mctString := fmt.Sprintf(MarketCloseTime, time.Now().Format("2006-01-02"))
	mct, err := time.ParseInLocation("2006-01-02 15:04:05", mctString, loc)
	if err != nil {
		return false, fmt.Errorf("error parsing market open time. %+v", err)
	}

	currentTime := time.Now()
	if currentTime.After(mot) && currentTime.Before(mct) && currentTime.Weekday() != 6 && currentTime.Weekday() != 7 {
		return true, nil
	}
	return false, nil

}

//IsWithInActualMarketOpenTime tells whether current time is within actual market open time.
func IsWithInActualMarketOpenTime() (bool, error) {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	amotString := fmt.Sprintf(ActualMarketOpenTime, time.Now().Format("2006-01-02"))
	amot, err := time.ParseInLocation("2006-01-02 15:04:05", amotString, loc)
	if err != nil {
		return false, fmt.Errorf("error parsing market open time. %+v", err)
	}

	mctString := fmt.Sprintf(MarketCloseTime, time.Now().Format("2006-01-02"))
	mct, err := time.ParseInLocation("2006-01-02 15:04:05", mctString, loc)
	if err != nil {
		return false, fmt.Errorf("error parsing market open time. %+v", err)
	}

	currentTime := time.Now()
	if currentTime.After(amot) && currentTime.Before(mct) && currentTime.Weekday() != 6 && currentTime.Weekday() != 0 {
		return true, nil
	}
	return false, nil

}

func IsWithInTimeRange(time1, time2 string) (bool, error) {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	t1String := fmt.Sprintf(time1, time.Now().Format("2006-01-02"))
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", t1String, loc)
	if err != nil {
		return false, fmt.Errorf("error parsing market open time. %+v", err)
	}

	t2String := fmt.Sprintf(time2, time.Now().Format("2006-01-02"))
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", t2String, loc)
	if err != nil {
		return false, fmt.Errorf("error parsing market open time. %+v", err)
	}

	currentTime := time.Now()
	if currentTime.After(t1) && currentTime.Before(t2) && int(currentTime.Weekday()) != 6 && int(currentTime.Weekday()) != 0 {
		return true, nil
	}
	return false, nil

}

func parseTime(format string, tstring string) (time.Time, error) {
	parsedTime, err := time.Parse(format, tstring)
	if err != nil {
		log.Fatalf("Error parsing market CloseTime Time: %+v", err)
		return time.Time{}, err
	}

	return parsedTime, nil
}

//GetDate returns the current date
func GetDate() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02")
}

//GetInterval returns interval in integer
func GetInterval(i string) int {
	switch i {
	case "5m":
		return 5

	case "10m":
		return 10

	case "15m":
		return 15

	default:
		return 0
	}

}

//GetSubscriptions returns the list of subscription IDs
func GetSubscriptions(stock string) []uint32 {
	token := []uint32{}

	stocks := strings.Split(stock, ",")

	for _, s := range stocks {
		sSlice := strings.Split(s, ";")
		if len(sSlice) >= 4 {
			token = append(token, getUnit32(sSlice[2]))
		}
	}

	return token
}

//GetSubscriptionsString returns '-' separated string of subscription IDs
func GetSubscriptionsString(stock string) string {
	token := ""

	stocks := strings.Split(stock, ",")

	for _, s := range stocks {
		sSlice := strings.Split(s, ";")
		if len(sSlice) >= 4 {
			token = token + "-" + sSlice[2]
		}
	}

	return token
}

//GetStocks returns formated 2-D array of stocks
func GetStocks(stock string) [][]string {
	result := [][]string{}
	stocks := strings.Split(stock, ",")

	for _, s := range stocks {
		sSlice := strings.Split(s, ";")
		result = append(result, sSlice)
	}

	return result
}

func getUnit32(str string) uint32 {
	u, _ := strconv.ParseUint(str, 10, 32)
	return uint32(u)
}
