package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harryoh/crypto-collector/bithumb"
	"github.com/harryoh/crypto-collector/bybit"
	"github.com/harryoh/crypto-collector/currency"
	"github.com/harryoh/crypto-collector/upbit"
	"github.com/joho/godotenv"
	"github.com/muesli/cache2go"
)

// Price :
type Price struct {
	Name      string
	Symbol    string
	Price     float64
	Timestamp int64
}

// Prices :
type Prices struct {
	UpbitPrice   Price
	BithumbPrice Price
	BybitPrice   Price
	UsdKrw       Price
	CreatedAt    int64
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
			time.Sleep(60)
			continue
		}
		val.Price = upbitTicker[0].TradePrice
		val.Timestamp = upbitTicker[0].Timestamp / 1000
		c <- *val
		time.Sleep(sleep)
	}
}

func bithumbLastPrice(sleep time.Duration, c chan Price) {
	val := &Price{
		Name:   "bithumb",
		Symbol: "BTC_KRW",
	}
	for {
		bithumbClient := bithumb.NewClient()
		bithumbTxHistory, err := bithumbClient.LastPrice(val.Symbol)
		if err != nil {
			fmt.Print(err)
			return
		}
		price, _ := strconv.ParseFloat(bithumbTxHistory.Data[0].Price, 64)
		val.Price = price
		val.Timestamp = time.Now().Unix()

		c <- *val
		time.Sleep(sleep)
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
		timestamp, _ := strconv.ParseFloat(bybitTicker.TimeNow, 64)
		val.Price = price
		val.Timestamp = int64(timestamp)

		c <- *val
		time.Sleep(sleep)
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
		time.Sleep(sleep)
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

// corsMiddleware :
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-control-Allow-Methods", "GET")
		c.Next()
	}
}

func lastPrice(c *gin.Context) {
	cache := _cache()
	name := c.Param("name")

	data, err := cache.Value(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, data.Data())
}

func allPrices(c *gin.Context) {
	var data *cache2go.CacheItem
	var err error
	cache := _cache()
	prices := &Prices{}

	data, err = cache.Value("upbit")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	prices.UpbitPrice = data.Data().(Price)

	data, err = cache.Value("bithumb")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	prices.BithumbPrice = data.Data().(Price)

	data, err = cache.Value("bybit")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	prices.BybitPrice = data.Data().(Price)

	data, err = cache.Value("usdkrw")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	prices.UsdKrw = data.Data().(Price)
	prices.CreatedAt = time.Now().Unix()

	c.JSON(http.StatusOK, prices)
}

func setPeriod(period map[string]time.Duration) {
	err := godotenv.Load()
	if err != nil {
		period["upbit"] = 5 * time.Second
		period["bithumb"] = 6 * time.Second
		period["bybit"] = 4 * time.Second
		period["usdkrw"] = 600 * time.Second
	} else {
		upbitPeriod, _ := strconv.Atoi(os.Getenv("UpbitPeriodSeconds"))
		period["upbit"] = time.Duration(upbitPeriod) * time.Second
		bithumbPeriod, _ := strconv.Atoi(os.Getenv("BithumbPeriodSeconds"))
		period["bithumb"] = time.Duration(bithumbPeriod) * time.Second
		bybitPeriod, _ := strconv.Atoi(os.Getenv("BybitPeriodSeconds"))
		period["bybit"] = time.Duration(bybitPeriod) * time.Second
		usdkrwPeriod, _ := strconv.Atoi(os.Getenv("UsdKrwPeriodSeconds"))
		period["usdkrw"] = time.Duration(usdkrwPeriod) * time.Second
	}
}

func main() {
	cache := _cache()

	period := make(map[string]time.Duration)
	setPeriod(period)
	go func() {
		ch := make(chan Price)

		// go upbitLastPrice(period["upbit"], ch)
		go upbitLastPrice(period["upbit"], ch)
		go bithumbLastPrice(period["bithumb"], ch)
		go bybitLastPrice(period["bybit"], ch)
		go usdPrice(period["usdkrw"], ch)

		for {
			select {
			case msg := <-ch:
				cache.Add(msg.Name, period[msg.Name]*2*time.Second, msg)
			}
		}
	}()

	router := gin.Default()
	router.Use(corsMiddleware())
	router.GET("/api/prices/:name", lastPrice)
	router.GET("/api/prices", allPrices)
	router.Run()

	log.Fatal(http.ListenAndServe(":8080", router))
}
