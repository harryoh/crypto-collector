package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/harryoh/crypto-collector/exchange/bithumb"
	"github.com/harryoh/crypto-collector/exchange/bybit"
	"github.com/harryoh/crypto-collector/exchange/currency"
	"github.com/harryoh/crypto-collector/exchange/upbit"
	"github.com/harryoh/crypto-collector/util"
	"github.com/joho/godotenv"
	"github.com/muesli/cache2go"
)

// Price :
type Price struct {
	Symbol    string
	Price     float64
	Timestamp int64
}

// Prices :
type Prices struct {
	Name      string
	Price     []Price
	Timestamp int64
}

// TotalPrices :
type TotalPrices struct {
	Currency     Prices
	BybitPrice   Prices
	UpbitPrice   Prices
	BithumbPrice Prices
	CreatedAt    int64
}

func upbitLastPrice(sleep time.Duration, c chan Prices) {
	val := &Prices{
		Name: "upbit",
	}
	for {
		val.Price = make([]Price, 0)
		markets := []string{"KRW-BTC", "KRW-ETH", "KRW-XRP"}

		upbitClient := upbit.NewClient()
		for _, market := range markets {
			upbitTicker, err := upbitClient.LastPrice(market)
			if err != nil {
				time.Sleep(sleep)
				continue
			}

			price := &Price{
				Symbol:    util.SymbolName(&market),
				Price:     upbitTicker[0].TradePrice,
				Timestamp: upbitTicker[0].Timestamp / 1000,
			}
			val.Price = append(val.Price, *price)
		}

		val.Timestamp = time.Now().Unix()

		c <- *val
		time.Sleep(sleep)
	}
}

func bithumbLastPrice(sleep time.Duration, c chan Prices) {
	val := &Prices{
		Name: "bithumb",
	}
	for {
		val.Price = make([]Price, 0)
		markets := []string{"BTC_KRW", "ETH_KRW", "XRP_KRW"}

		loc, _ := time.LoadLocation("Asia/Seoul")
		bithumbClient := bithumb.NewClient()
		for _, market := range markets {
			bithumbTxHistory, err := bithumbClient.LastPrice(market)
			if err != nil {
				fmt.Print(err)
				return
			}

			lastPrice, _ := strconv.ParseFloat(bithumbTxHistory.Data[0].Price, 64)
			kst, _ := time.ParseInLocation("2006-01-02 15:04:05", bithumbTxHistory.Data[0].TransactionDate, loc)
			price := &Price{
				Symbol:    util.SymbolName(&market),
				Price:     lastPrice,
				Timestamp: kst.Unix(),
			}

			val.Price = append(val.Price, *price)
		}

		val.Timestamp = time.Now().Unix()

		c <- *val
		time.Sleep(sleep)
	}
}

func bybitLastPrice(sleep time.Duration, c chan Prices) {
	val := &Prices{
		Name: "bybit",
	}
	for {
		val.Price = make([]Price, 0)

		bybitClient := bybit.NewClient()
		bybitTicker, err := bybitClient.LastPrice("")
		if err != nil {
			fmt.Print(err)
			return
		}

		for _, result := range bybitTicker.Result {
			symbol := util.SymbolName(&result.Symbol)
			if symbol == "" {
				continue
			}

			lastPrice, _ := strconv.ParseFloat(result.LastPrice, 64)
			timestamp, _ := strconv.ParseFloat(bybitTicker.TimeNow, 64)

			price := &Price{
				Symbol:    util.SymbolName(&result.Symbol),
				Price:     lastPrice,
				Timestamp: int64(timestamp),
			}

			val.Price = append(val.Price, *price)
		}

		val.Timestamp = time.Now().Unix()

		c <- *val
		time.Sleep(sleep)
	}
}

func currencyRate(sleep time.Duration, c chan Prices) {
	val := &Prices{
		Name: "currency",
	}
	for {
		val.Price = make([]Price, 0)
		markets := []string{"USDKRW"}

		currencyClient := currency.NewClient()
		for _, market := range markets {
			rate, err := currencyClient.CurrencyRate(market)
			if err != nil {
				return
			}

			price := &Price{
				Symbol:    market,
				Price:     rate.USDKRW[0],
				Timestamp: rate.Update / 1000,
			}
			val.Price = append(val.Price, *price)
		}

		val.Timestamp = time.Now().Unix()
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
	totalPrices := &TotalPrices{}

	data, err = cache.Value("upbit")
	if err == nil {
		totalPrices.UpbitPrice = data.Data().(Prices)
	}

	data, err = cache.Value("bithumb")
	if err == nil {
		totalPrices.BithumbPrice = data.Data().(Prices)
	}

	data, err = cache.Value("bybit")
	if err == nil {
		totalPrices.BybitPrice = data.Data().(Prices)
	}

	data, err = cache.Value("currency")
	if err == nil {
		totalPrices.Currency = data.Data().(Prices)
	}

	totalPrices.CreatedAt = time.Now().Unix()

	c.JSON(http.StatusOK, totalPrices)
}

func setPeriod(period map[string]time.Duration) {
	err := godotenv.Load()
	if err != nil {
		period["upbit"] = 5 * time.Second
		period["bithumb"] = 6 * time.Second
		period["bybit"] = 4 * time.Second
		period["currency"] = 600 * time.Second
	} else {
		upbitPeriod, _ := strconv.Atoi(os.Getenv("UpbitPeriodSeconds"))
		period["upbit"] = time.Duration(upbitPeriod) * time.Second
		bithumbPeriod, _ := strconv.Atoi(os.Getenv("BithumbPeriodSeconds"))
		period["bithumb"] = time.Duration(bithumbPeriod) * time.Second
		bybitPeriod, _ := strconv.Atoi(os.Getenv("BybitPeriodSeconds"))
		period["bybit"] = time.Duration(bybitPeriod) * time.Second
		currencyPeriod, _ := strconv.Atoi(os.Getenv("CurrencyPeriodSeconds"))
		period["currency"] = time.Duration(currencyPeriod) * time.Second
	}
}

func main() {
	cache := _cache()

	period := make(map[string]time.Duration)
	setPeriod(period)
	go func() {
		ch := make(chan Prices)

		go upbitLastPrice(period["upbit"], ch)
		go bithumbLastPrice(period["bithumb"], ch)
		go bybitLastPrice(period["bybit"], ch)
		go currencyRate(period["currency"], ch)

		for {
			select {
			case msg := <-ch:
				cache.Add(msg.Name, period[msg.Name]*2*time.Second, msg)
			}
		}
	}()

	router := gin.Default()
	router.Use(corsMiddleware())
	router.Use(static.Serve("/", static.LocalFile("./ui/build", true)))
	router.GET("/api/prices/:name", lastPrice)
	router.GET("/api/prices", allPrices)
	router.Run()

	log.Fatal(http.ListenAndServe(":8080", router))
}
