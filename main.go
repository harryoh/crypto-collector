package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/harryoh/crypto-collector/bybit"
	"github.com/harryoh/crypto-collector/currency"
	"github.com/harryoh/crypto-collector/upbit"
)

func main() {
	loc, _ := time.LoadLocation("Asia/Seoul")

	upbitClient := upbit.NewClient()
	upbitTicker, err := upbitClient.LastPrice("KRW-BTC")
	if err != nil {
		return
	}
	fmt.Printf("%f\n", upbitTicker[0].TradePrice)
	fmt.Println(time.Unix(upbitTicker[0].Timestamp/1000, 0).In(loc))

	bybitClient := bybit.NewClient()
	bybitTicker, err := bybitClient.Tickers("BTCUSD")
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println(bybitTicker.Result[0].LastPrice)
	// fmt.Printf("%T\n", bybitTicker.TimeNow)
	t, _ := strconv.ParseFloat(bybitTicker.TimeNow, 64)
	fmt.Println(time.Unix(int64(t), 0).In(loc))

	currencyClient := currency.NewClient()
	rate, err := currencyClient.ExchangeRate("USDKRW")
	if err != nil {
		return
	}
	fmt.Printf("%f\n", rate.USDKRW[0])
	fmt.Println(time.Unix(rate.Update/1000, 0).In(loc))
}
