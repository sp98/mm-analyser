package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/sp98/analyzer/pkg/analyze"

	"github.com/sp98/analyzer/pkg/utility"
)

const (
	//StocksEnv is the enviornment variable to get all the stocks.
	StocksEnv = "STOCKS"
	//IntervalEnv is the environment variable for 5 minutes interval
	IntervalEnv = "INTERVAL_5M"
)

var (
	interval string
	stocks   string
)

func init() {
	interval = os.Getenv(IntervalEnv)
	if interval == "" {
		log.Fatalf("error reading env variables %q", IntervalEnv)
		panic(1)
	}
	stocks = os.Getenv(StocksEnv)
	if stocks == "" {
		log.Fatalf("error reading env variables %q", StocksEnv)
		panic(1)
	}
}

func main() {
	log.Println("--- START ANALYSER  ---")
	isMarketOpen()
	setup()
}

func setup() {
	//Get all the stocks in a 2D array with format - Instrument Name, Sybmol, Token, Exchange, Interval

	intervalInt := utility.GetInterval("3m")
	waitFor := utility.WaitBeforeAnalysis(intervalInt)
	if waitFor > 0 {
		log.Printf("Wait for %.2f minutes before starting", float64(waitFor)/60)
		time.Sleep(time.Second * time.Duration(waitFor))
	}

	t := time.NewTicker(time.Minute * time.Duration(intervalInt))

	time.Sleep(time.Second * 3) //Adding a wait for the continuous query to run
	startAnalysis()

	log.Printf("Analysis Start Time: %+v <---> Analysis Stop Time: %+v", time.Now(), fmt.Sprintf(utility.MarketCloseTime, utility.GetDate()))
	for alive := true; alive; {
		stamp := <-t.C
		log.Printf("Starting Analysis at %+v", stamp.Format(utility.TstringFormat))
		time.Sleep(time.Second * 3) //Adding a wait for the continuous query to run
		startAnalysis()
	}

}

func startAnalysis() {
	var wg sync.WaitGroup
	stocks := utility.GetStocks(stocks)
	result := &analyze.Result{Mux: &sync.Mutex{}}

	for _, stock := range stocks {
		wg.Add(1)
		ohlc, err := analyze.GetOHLC(stock[2], interval)
		if err != nil {
			log.Printf("error in setup for stock %q. %+v", stock[0], err)
		}
		insturment := analyze.NewInsturment(stock[0], stock[1], stock[2], stock[3], ohlc)

		go insturment.Analyze(result, &wg)
	}

	wg.Wait() //wait for all the go routines to finish

	log.Printf("Result - %+v", result)
	isMarketOpen()

}

func isMarketOpen() {
	for {
		t, _ := utility.IsWithInMarketOpenTime()
		if t {
			log.Println(" withhin market open time.")
			break
		}
		log.Println("not withhin market open time. Waiting...")
		time.Sleep(10 * time.Second)
	}

}
