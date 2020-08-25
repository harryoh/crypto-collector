package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/harryoh/crypto-collector/bybit"
	"github.com/harryoh/crypto-collector/currency"
	"github.com/harryoh/crypto-collector/upbit"
)

// Price :
type Price struct {
	price     float64
	timestamp int64
}

func upbitLastPrice(sleep time.Duration, c chan Price) {
	val := &Price{}
	for {
		upbitClient := upbit.NewClient()
		upbitTicker, err := upbitClient.LastPrice("KRW-BTC")
		if err != nil {
			time.Sleep(60 * time.Second)
			continue
		}
		val.price = upbitTicker[0].TradePrice
		val.timestamp = upbitTicker[0].Timestamp
		c <- *val
		time.Sleep(sleep * time.Second)
	}
}

func bybitLastPrice(sleep time.Duration, c chan Price) {
	val := &Price{}
	for {
		bybitClient := bybit.NewClient()
		bybitTicker, err := bybitClient.LastPrice("BTCUSD")
		if err != nil {
			fmt.Print(err)
			return
		}
		price, _ := strconv.ParseFloat(bybitTicker.Result[0].LastPrice, 64)
		timestamp, _ := strconv.Atoi(bybitTicker.TimeNow)
		val.price = price
		val.timestamp = int64(timestamp)

		c <- *val
		time.Sleep(sleep * time.Second)
	}
}

func usdPrice(sleep time.Duration, c chan Price) {
	val := &Price{}
	for {
		currencyClient := currency.NewClient()
		rate, err := currencyClient.ExchangeRate("USDKRW")
		if err != nil {
			return
		}
		val.price = rate.USDKRW[0]
		val.timestamp = rate.Update / 1000
		c <- *val
		time.Sleep(sleep * time.Second)
	}
}

func main() {
	upbitCh := make(chan Price)
	bybitCh := make(chan Price)
	currencyCh := make(chan Price)

	go upbitLastPrice(1, upbitCh)
	go bybitLastPrice(2, bybitCh)
	go usdPrice(3, currencyCh)

	for {
		select {
		case upbitMsg := <-upbitCh:
			fmt.Printf("Upbit: %f\n", upbitMsg.price)
		case bybitMsg := <-bybitCh:
			fmt.Printf("Bybit: %f\n", bybitMsg.price)
		case currencyMsg := <-currencyCh:
			fmt.Printf("USDKRW: %f %v\n", currencyMsg.price, currencyMsg.timestamp)
		}
	}
}
