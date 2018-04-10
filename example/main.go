package main

import (
	"fmt"
	"time"
	"github.com/cjongseok/slog"
	"github.com/cjongseok/fetch-bittrex"
	"flag"
)

var (
	apikey = flag.String("apikey", "", "Bittrex API Key")
	apisec = flag.String("apisec", "", "Bittrex API Secrets")
)

func main() {
	flag.Parse()
	if *apikey == "" || *apisec == "" {
		fmt.Println("ERROR: -apikey and -apisec are required.")
		flag.PrintDefaults()
		return
	}

	fmt.Println("Start")
	fetcher.Start(*apikey, *apisec)
	fmt.Printf("%s: Fetched? %v\n", time.Now(), fetcher.Fetched())
	fetcher.WaitForFetching()
	fmt.Printf("%s: Fetched? %v\n", time.Now(), fetcher.Fetched())
	fmt.Println("Coins:", slog.Stringify(fetcher.All()))
	fmt.Println("Size:", fetcher.Size())
	fmt.Println("Get BTC-BAT:", slog.Stringify(fetcher.Get("BTC-BAT")))
	fmt.Println("Now:", time.Now())
	fmt.Println("Next coin fetching time:", fetcher.NextCoinFetchTime())
	fetcher.Close()
	fmt.Println("Closed")
}



