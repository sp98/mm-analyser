package main

import (
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
)

func main() {
	log.Println("Analyszer")
	setup()
}

func setup() {
	//Get all the stocks in a 2D array with format - Instrument Name, Sybmol, Token, Exchange, Interval
	stocks := utility.GetStocks(os.Getenv(StocksEnv))
	result := &analyze.Result{Mux: &sync.Mutex{}}

	for _, stock := range stocks {
		ohlc, err := analyze.GetOHLC(stock[2], "5m")
		if err != nil {
			log.Printf("error fetch ohlc. %+v", err)
		}
		insturment := analyze.NewInsturment(stock[0], stock[1], stock[2], stock[3], ohlc)

		go insturment.Analyze(result)
	}

	//wait for all the go routines to finish
	time.Sleep(5 * time.Second)
	log.Printf("Result - %+v", result)

}
