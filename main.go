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
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
	Symbol               string
	Price                float64
	FundingRate          float64
	PredictedFundingRate float64
	Timestamp            int64
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
	Rule         rule
	CreatedAt    int64
}

type telegramKey struct {
	ChatID int64
	Token  string
}

type rule struct {
	Use      bool
	AlarmMax float64
	AlarmMin float64
}

type envs struct {
	Period         map[string]time.Duration
	Monitor        telegramKey
	Alarm          telegramKey
	CurrencyAPIKey string
	Rule           rule
}

type sendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func upbitLastPrice(env *envs, c chan Prices) {
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
				fmt.Print("upbitLastPrice Error: ")
				fmt.Println(err)
				time.Sleep(env.Period["upbit"])
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
		time.Sleep(env.Period["upbit"])
	}
}

func bithumbLastPrice(env *envs, c chan Prices) {
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
				fmt.Print("bithumbLastPrice Error: ")
				fmt.Println(err)
				time.Sleep(env.Period["bithumb"])
				continue
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
		time.Sleep(env.Period["bithumb"])
	}
}

func bybitLastPrice(env *envs, c chan Prices) {
	val := &Prices{
		Name: "bybit",
	}
	for {
		val.Price = make([]Price, 0)

		bybitClient := bybit.NewClient()
		bybitTicker, err := bybitClient.LastPrice("")
		if err != nil {
			fmt.Print("bybitLastPrice Error: ")
			fmt.Println(err)
			time.Sleep(env.Period["bybit"])
			continue
		}

		for _, result := range bybitTicker.Result {
			symbol := util.SymbolName(&result.Symbol)
			if symbol == "" {
				continue
			}

			lastPrice, _ := strconv.ParseFloat(result.LastPrice, 64)
			timestamp, _ := strconv.ParseFloat(bybitTicker.TimeNow, 64)
			fundingrate, _ := strconv.ParseFloat(result.FundingRate, 64)
			predictedFundingRate, _ := strconv.ParseFloat(result.PredictedFundingRate, 64)

			price := &Price{
				Symbol:               util.SymbolName(&result.Symbol),
				Price:                lastPrice,
				Timestamp:            int64(timestamp),
				FundingRate:          fundingrate,
				PredictedFundingRate: predictedFundingRate,
			}

			val.Price = append(val.Price, *price)
		}

		val.Timestamp = time.Now().Unix()

		c <- *val
		time.Sleep(env.Period["bybit"])
	}
}

func currencyRate(env *envs, c chan Prices) {
	val := &Prices{
		Name: "currency",
	}
	for {
		val.Price = make([]Price, 0)
		markets := []string{"USD_KRW"}

		currencyClient := currency.NewClient()
		for _, market := range markets {
			rate, err := currencyClient.CurrencyRate(market, env.CurrencyAPIKey)
			if err != nil {
				fmt.Print("currencyRate Error: ")
				fmt.Println(err)
				time.Sleep(env.Period["currency"])
				continue
			}

			price := &Price{
				Symbol:    market,
				Price:     rate.USDKRW,
				Timestamp: time.Now().Unix(),
			}
			val.Price = append(val.Price, *price)
		}
		val.Timestamp = time.Now().Unix()
		c <- *val
		time.Sleep(env.Period["currency"])
	}
}

func premiumRate(bybit float64, desc float64) float64 {
	return (desc - bybit*1200) / desc * 100
}

func sendMonitorMessage(env *envs) {
	if env.Alarm.Token == "" || env.Alarm.ChatID == 0 {
		log.Println("Key is invalid for alarm")
		return
	}

	alarmBot, err := tgbotapi.NewBotAPI(env.Alarm.Token)
	if err != nil {
		panic(err)
	}

	cache := _cache()
	alarmBot.Debug = false

	cnt := 0
	for {
		time.Sleep(env.Period["alarm"])

		var data *cache2go.CacheItem
		data, err = cache.Value("rule")
		if err != nil {
			continue
		}
		env.Rule = data.Data().(rule)

		if env.Rule.Use != true {
			continue
		}

		totalPrices := readPrices()

		if len(totalPrices.BybitPrice.Price) < 1 {
			fmt.Println("Bybit Prices is NULL!")
			continue
		}

		info := "http://home.5004.pe.kr:8080\n" +
			"KRWUSD:" + strconv.FormatFloat(totalPrices.Currency.Price[0].Price, 'f', -1, 64) +
			" FixKRWUSD: 1200\n"

		content := "[Bybit] " +
			" BTC:" + strconv.FormatFloat(totalPrices.BybitPrice.Price[0].Price, 'f', -1, 64) +
			"(" + strconv.FormatFloat(totalPrices.BybitPrice.Price[0].FundingRate, 'f', -1, 64) + ")" +
			" ETH:" + strconv.FormatFloat(totalPrices.BybitPrice.Price[1].Price, 'f', -1, 64) +
			"(" + strconv.FormatFloat(totalPrices.BybitPrice.Price[0].FundingRate, 'f', -1, 64) + ")"

		var premiumRateBithumbBTC float64
		var premiumRateBithumbETH float64
		if len(totalPrices.BithumbPrice.Price) > 1 {
			premiumRateBithumbBTC = premiumRate(totalPrices.BybitPrice.Price[0].Price, totalPrices.BithumbPrice.Price[0].Price)
			premiumRateBithumbETH = premiumRate(totalPrices.BybitPrice.Price[1].Price, totalPrices.BithumbPrice.Price[1].Price)
			content += "\n[Bithumb] " +
				" BTC:" + strconv.FormatFloat(totalPrices.BithumbPrice.Price[0].Price, 'f', -1, 64) +
				"(" + strconv.FormatFloat(premiumRateBithumbBTC, 'f', 3, 64) + "%)" +
				" ETH:" + strconv.FormatFloat(totalPrices.BybitPrice.Price[1].Price, 'f', -1, 64) +
				"(" + strconv.FormatFloat(premiumRate(totalPrices.BybitPrice.Price[1].Price, totalPrices.BithumbPrice.Price[1].Price), 'f', 3, 64) + "%)"
		}
		var premiumRateUpbitBTC float64
		var premiumRateUpbitETH float64
		if len(totalPrices.UpbitPrice.Price) > 1 {
			premiumRateUpbitBTC = premiumRate(totalPrices.BybitPrice.Price[0].Price, totalPrices.UpbitPrice.Price[0].Price)
			premiumRateUpbitETH = premiumRate(totalPrices.BybitPrice.Price[1].Price, totalPrices.UpbitPrice.Price[1].Price)
			content += "\n[Upbit] " +
				" BTC:" + strconv.FormatFloat(totalPrices.UpbitPrice.Price[0].Price, 'f', -1, 64) +
				"(" + strconv.FormatFloat(premiumRateUpbitBTC, 'f', 3, 64) + "%)" +
				" ETH:" + strconv.FormatFloat(totalPrices.BybitPrice.Price[1].Price, 'f', -1, 64) +
				"(" + strconv.FormatFloat(premiumRate(totalPrices.BybitPrice.Price[1].Price, totalPrices.UpbitPrice.Price[1].Price), 'f', 3, 64) + "%)"
		}

		if cnt%50 == 0 {
			content = info + content
		}

		// fmt.Println(env.Rule.Use, env.Rule.AlarmMin, env.Rule.AlarmMax, premiumRateUpbitBTC, premiumRateBithumbETH, premiumRateUpbitBTC, premiumRateUpbitETH)
		ruleText := "RULE [ Max:" + strconv.FormatFloat(env.Rule.AlarmMax, 'f', -1, 64) +
			" Min:" + strconv.FormatFloat(env.Rule.AlarmMin, 'f', -1, 64) + " ]\n"
		if len(totalPrices.BithumbPrice.Price) > 1 {
			if premiumRateBithumbBTC <= env.Rule.AlarmMin || premiumRateBithumbBTC >= env.Rule.AlarmMax ||
				premiumRateBithumbETH <= env.Rule.AlarmMin || premiumRateBithumbETH >= env.Rule.AlarmMax {
				msg := tgbotapi.NewMessage(env.Alarm.ChatID, ruleText+content)
				alarmBot.Send(msg)
				continue
			}
		}

		if len(totalPrices.UpbitPrice.Price) > 1 {
			if premiumRateUpbitBTC <= env.Rule.AlarmMin || premiumRateUpbitBTC >= env.Rule.AlarmMax ||
				premiumRateUpbitETH <= env.Rule.AlarmMin || premiumRateUpbitETH >= env.Rule.AlarmMax {
				msg := tgbotapi.NewMessage(env.Alarm.ChatID, ruleText+content)
				alarmBot.Send(msg)
				continue
			}
		}
		cnt++
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

func readPrices() (totalPrices *TotalPrices) {
	var data *cache2go.CacheItem
	var err error
	cache := _cache()
	totalPrices = &TotalPrices{}

	data, err = cache.Value("upbit")
	totalPrices.UpbitPrice.Name = "upbit"
	totalPrices.UpbitPrice.Price = make([]Price, 0)
	if err == nil {
		totalPrices.UpbitPrice = data.Data().(Prices)
	}

	data, err = cache.Value("bithumb")
	totalPrices.BithumbPrice.Name = "bithumb"
	totalPrices.BithumbPrice.Price = make([]Price, 0)
	if err == nil {
		totalPrices.BithumbPrice = data.Data().(Prices)
	}

	data, err = cache.Value("bybit")
	totalPrices.BybitPrice.Name = "bybit"
	totalPrices.BybitPrice.Price = make([]Price, 0)
	if err == nil {
		totalPrices.BybitPrice = data.Data().(Prices)
	}

	data, err = cache.Value("currency")
	totalPrices.Currency.Name = "currency"
	totalPrices.Currency.Price = make([]Price, 0)
	if err == nil {
		totalPrices.Currency = data.Data().(Prices)
	}

	data, err = cache.Value("rule")
	if err == nil {
		totalPrices.Rule = data.Data().(rule)
	}

	totalPrices.CreatedAt = time.Now().Unix()
	return
}

func allPrices(c *gin.Context) {
	totalPrices := readPrices()
	c.JSON(http.StatusOK, totalPrices)
}

func setRule(c *gin.Context) {
	var json rule
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cache := _cache()
	cache.Add("rule", 0, json)

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func setEnvs(env *envs) {
	godotenv.Load()
	// Default Value
	env.Period["upbit"] = 4 * time.Second
	env.Period["bithumb"] = 5 * time.Second
	env.Period["bybit"] = 3 * time.Second
	env.Period["currency"] = 60 * 60 * time.Second
	env.Period["alarm"] = 10 * time.Second

	upbitPeriod, _ := strconv.Atoi(os.Getenv("UpbitPeriodSeconds"))
	if upbitPeriod > 0 {
		env.Period["upbit"] = time.Duration(upbitPeriod) * time.Second
	}
	bithumbPeriod, _ := strconv.Atoi(os.Getenv("BithumbPeriodSeconds"))
	if bithumbPeriod > 0 {
		env.Period["bithumb"] = time.Duration(bithumbPeriod) * time.Second
	}
	bybitPeriod, _ := strconv.Atoi(os.Getenv("BybitPeriodSeconds"))
	if bybitPeriod > 0 {
		env.Period["bybit"] = time.Duration(bybitPeriod) * time.Second
	}
	currencyPeriod, _ := strconv.Atoi(os.Getenv("CurrencyPeriodSeconds"))
	if currencyPeriod > 0 {
		env.Period["currency"] = time.Duration(currencyPeriod) * time.Second
	}

	messagePeriod, _ := strconv.Atoi(os.Getenv("MessagePeriodSeconds"))
	if messagePeriod > 0 {
		env.Period["alarm"] = time.Duration(messagePeriod) * time.Second
	}

	env.Alarm.ChatID, _ = strconv.ParseInt(os.Getenv("AlarmChatID"), 10, 64)
	env.Alarm.Token = os.Getenv("AlarmToken")

	env.CurrencyAPIKey = os.Getenv("CurrencyAPIKey")

	env.Rule.Use, _ = strconv.ParseBool(os.Getenv("RuleAlarmUse"))
	env.Rule.AlarmMax, _ = strconv.ParseFloat(os.Getenv("RuleAlarmMax"), 64)
	env.Rule.AlarmMin, _ = strconv.ParseFloat(os.Getenv("RuleAlarmMin"), 64)

	cache := _cache()
	cache.Add("rule", 0, env.Rule)
}

func main() {
	cache := _cache()

	env := &envs{
		Period: make(map[string]time.Duration),
	}

	setEnvs(env)

	go func() {
		ch := make(chan Prices)

		go upbitLastPrice(env, ch)
		go bithumbLastPrice(env, ch)
		go bybitLastPrice(env, ch)
		go currencyRate(env, ch)
		go sendMonitorMessage(env)

		for {
			select {
			case msg := <-ch:
				cache.Add(msg.Name, env.Period[msg.Name]*20, msg)
			}
		}
	}()

	router := gin.Default()
	router.Use(corsMiddleware())
	router.Use(static.Serve("/", static.LocalFile("./ui/build", true)))
	router.GET("/api/prices/:name", lastPrice)
	router.GET("/api/prices", allPrices)
	router.POST("/api/rule", setRule)
	router.Run()

	log.Fatal(http.ListenAndServe(":8080", router))
}
