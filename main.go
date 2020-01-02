package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sp98/analyzer/pkg/analyze"
	"github.com/sp98/analyzer/pkg/utility"
)

const (
	//IntervalEnv is the environment variable for 5 minutes interval
	IntervalEnv = "ANALYZER_INTERVAL"
	//StocksEnv is the environment variable for list of stocks to analyse
	StocksEnv = "STOCKS"
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
	log.Println("command line args - ", os.Args[1:])
	if len(os.Args) > 1 {
		analyze.FromTime = os.Args[1]
		analyze.ToTime = os.Args[2]
	}

	log.Println("--- START ANALYSER  ---")
	utility.IsMarketOpen()
	setup()
}

func setup() {
	//Get all the stocks in a 2D array with format - Instrument Name, Sybmol, Token, Exchange, Interval
	intervalInt := utility.GetInterval(interval)
	waitFor := utility.WaitBeforeAnalysis(intervalInt)
	if waitFor > 0 {
		log.Printf("Wait for %.2f minutes before starting", float64(waitFor)/60)
		time.Sleep(time.Second * time.Duration(waitFor))
	}

	t := time.NewTicker(time.Minute * time.Duration(intervalInt))

	time.Sleep(time.Second * 3) //Adding a wait for the continuous query to run
	analyze.StartAnalysis(stocks, interval)

	log.Printf("Analysis Start Time: %+v <---> Analysis Stop Time: %+v", time.Now(), fmt.Sprintf(utility.MarketCloseTime, utility.GetDate()))
	for alive := true; alive; {
		stamp := <-t.C
		log.Printf("Starting Analysis at %+v", stamp.Format(utility.TstringFormat))
		time.Sleep(time.Second * 3) //Adding a wait for the continuous query to run
		analyze.StartAnalysis(stocks, interval)
	}

}
