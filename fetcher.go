package fetcher

import (
	"sync"
	"time"
	"github.com/cjongseok/slog"
	api "github.com/toorop/go-bittrex"
)

const (
	logTag = "[fetch-bittrex]"
	defaultDelay = 30 * time.Second
	retryDelay   = 5 * time.Second
)

var bclient *api.Bittrex
var normalDelay time.Duration
var m sync.Mutex
var wg sync.WaitGroup
var coinFetchWg sync.WaitGroup
var started bool
var coinFetched bool
var coinFetchInterrupt chan struct{}
var all map[string]api.MarketSummary
var coinFetchedTime time.Time
var nextCoinFetchTime time.Time

func Start(apiKey, apiSecret string) chan []api.MarketSummary{
	return StartLimit(apiKey, apiSecret, defaultDelay)
}
func StartLimit(apiKey, apiSecret string, delay time.Duration) chan []api.MarketSummary{
	m.Lock()
	defer m.Unlock()
	if started {
		return nil
	}
	bclient = api.New(apiKey, apiSecret)
	normalDelay = delay
	coinFetchInterrupt = make(chan struct{})
	all = make(map[string]api.MarketSummary)
	wg = sync.WaitGroup{}
	coinFetchWg = sync.WaitGroup{}
	wg.Add(1)
	coinFetchWg.Add(1)
	coinUpdates := fetchCoin()
	started = true
	return coinUpdates
}
func WaitForFetching() {
	if !started {
		return
	}
	coinFetchWg.Wait()
}
func fetchCoin() chan []api.MarketSummary {
	out := make(chan []api.MarketSummary)
	go func() {
		defer wg.Done()
		var fetchDelay time.Duration
		streaming := true
		streamMutex := sync.Mutex{}
		isStreaming := func() bool {
			streamMutex.Lock()
			defer streamMutex.Unlock()
			return streaming
		}
		stopStreaming := func() {
			streamMutex.Lock()
			defer streamMutex.Unlock()
			streaming = false
		}
		fetch := func() (changed []api.MarketSummary) {
			coins, err := bclient.GetMarketSummaries()
			if err != nil {
				slog.Logf(logTag, "coin fetch failure: %s\n", err)
				fetchDelay = retryDelay
				return
			}
			for _, coin := range coins {
				old, ok := all[coin.MarketName]
				if !ok || (ok && old != coin) {
					all[coin.MarketName] = coin
					changed = append(changed, coin)
				}
			}
			if !coinFetched {
				coinFetched = true
				coinFetchWg.Done()
			}
			fetchDelay = normalDelay
			coinFetchedTime = time.Now()
			return
		}
		stream := func(coins []api.MarketSummary, to chan []api.MarketSummary) {
			if len(coins) < 1 {
				return
			}
			defer func() {
				// handle panics on pushing to closed channel
				recover()
			}()
			if isStreaming() {
				to <- coins
			}
		}
		newCoins := fetch()
		go stream(newCoins, out)
		for {
			nextCoinFetchTime = coinFetchedTime.Add(fetchDelay)
			select {
			case <-coinFetchInterrupt:
				stopStreaming()
				close(out)
				return
			case <-time.After(fetchDelay):
				updated := fetch()
				go stream(updated, out)
			}
		}
	}()
	return out
}
func Fetched() bool {
	return coinFetched
}
func Size() int {
	if coinFetched {
		return len(all)
	}
	return 0
}
func All() map[string]api.MarketSummary {
	if coinFetched {
		return all
	}
	return nil
}
func Get(market string) api.MarketSummary {
	if coinFetched {
		return all[market]
	}
	return api.MarketSummary{}
}
func NextCoinFetchTime() time.Time {
	return nextCoinFetchTime
}
func Close() {
	close(coinFetchInterrupt)
	wg.Wait()
}
