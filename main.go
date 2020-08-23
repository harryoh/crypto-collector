package main

import (
	"fmt"

	"github.com/harryoh/crypto-collector/upbit"
)

func main() {
	upbitClient := upbit.NewClient()

	markets, err := upbitClient.Markets()
	if err != nil {
		return
	}
	fmt.Println(markets[0].Market) // KRW-BTC

	minCandle, err := upbitClient.MinuteCandles(1, "KRW-BTC")
	if err != nil {
		return
	}
	fmt.Printf("%+v\n", minCandle[0].TradePrice)
	fmt.Println(minCandle[0].Market)
}
