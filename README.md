coinfetcher
===
Cache for Bittrex ticks.
It periodically fetch coin ticks from [Bittrex](https://bittrex.com) using [toorop/go-bittrex](https://github.com/toorop/go-bittrex).

Usage
---
Turn on the fetcher
```go
fetcher.Start(apiKey, apiSecerets)  // default fetching delay is 30 seconds.
fetcher.WaitForFetching()
```
And get ticks.
```go
fetcher.Get("BTC-LTC")  // get recent BTC-LTC tick
fetcher.All()           // all the recent ticks
```
