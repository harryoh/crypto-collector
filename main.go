package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/harryoh/crypto-collector/bybit"
	"github.com/harryoh/crypto-collector/currency"
	"github.com/harryoh/crypto-collector/upbit"
	"github.com/muesli/cache2go"
)

// Price :
type Price struct {
	Name      string
	Symbol    string
	Price     float64
	Timestamp int64
}

func upbitLastPrice(sleep time.Duration, c chan Price) {
	val := &Price{
		Name:   "upbit",
		Symbol: "KRW-BTC",
	}
	for {
		upbitClient := upbit.NewClient()
		upbitTicker, err := upbitClient.LastPrice(val.Symbol)
		if err != nil {
			time.Sleep(60 * time.Second)
			continue
		}
		val.Price = upbitTicker[0].TradePrice
		val.Timestamp = upbitTicker[0].Timestamp
		c <- *val
		time.Sleep(sleep * time.Second)
	}
}

func bybitLastPrice(sleep time.Duration, c chan Price) {
	val := &Price{
		Name:   "bybit",
		Symbol: "BTCUSD",
	}
	for {
		bybitClient := bybit.NewClient()
		bybitTicker, err := bybitClient.LastPrice(val.Symbol)
		if err != nil {
			fmt.Print(err)
			return
		}
		price, _ := strconv.ParseFloat(bybitTicker.Result[0].LastPrice, 64)
		timestamp, _ := strconv.Atoi(bybitTicker.TimeNow)
		val.Price = price
		val.Timestamp = int64(timestamp)

		c <- *val
		time.Sleep(sleep * time.Second)
	}
}

func usdPrice(sleep time.Duration, c chan Price) {
	val := &Price{
		Name:   "usdkrw",
		Symbol: "USDKRW",
	}
	for {
		currencyClient := currency.NewClient()
		rate, err := currencyClient.ExchangeRate(val.Symbol)
		if err != nil {
			return
		}
		val.Price = rate.USDKRW[0]
		val.Timestamp = rate.Update / 1000
		c <- *val
		time.Sleep(sleep * time.Second)
	}
}

func _cacheValue(key string) string {
	cache := _cache()
	data, err := cache.Value(key)
	if err != nil {
		fmt.Println(err)
	}
	res, err := json.Marshal(data.Data())
	if err != nil {
		fmt.Println(err)
	}

	return string(res)

}

func _cache() *cache2go.CacheTable {
	return cache2go.Cache("crypto")
}

func lastPrice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cache := _cache()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	data, err := cache.Value(vars["name"])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	res, err := json.Marshal(data.Data())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(res)
}

func main() {
	cache := _cache()
	period := map[string]time.Duration{
		"upbit":  6,
		"bybit":  5,
		"usdkrw": 60,
	}
	go func() {
		ch := make(chan Price)

		go upbitLastPrice(period["upbit"], ch)
		go bybitLastPrice(period["bybit"], ch)
		go usdPrice(period["usdkrw"], ch)

		for {
			select {
			case msg := <-ch:
				cache.Add(msg.Name, period[msg.Name]*2*time.Second, msg)
			}
		}
	}()

	router := mux.NewRouter()
	router.Use(mux.CORSMethodMiddleware(router))
	router.HandleFunc("/api/lastprice/{name}", lastPrice)
	log.Fatal(http.ListenAndServe(":8080", router))
}
