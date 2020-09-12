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
	Symbol      string
	Price       float64
	FundingRate float64
	Timestamp   int64
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

type telegramKey struct {
	ChatID int64
	Token  string
}

type envs struct {
	Period  map[string]time.Duration
	Monitor telegramKey
	Alarm   telegramKey
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
				fmt.Print("Error: ")
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
				fmt.Print("Error: ")
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
			fmt.Print("Error: ")
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

			price := &Price{
				Symbol:      util.SymbolName(&result.Symbol),
				Price:       lastPrice,
				Timestamp:   int64(timestamp),
				FundingRate: fundingrate,
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
		markets := []string{"USDKRW"}

		currencyClient := currency.NewClient()
		for _, market := range markets {
			rate, err := currencyClient.CurrencyRate(market)
			if err != nil {
				fmt.Print("Error: ")
				fmt.Println(err)
				time.Sleep(env.Period["currency"])
				continue
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
		time.Sleep(env.Period["currency"])
	}
}

func premiumRate(bybit float64, desc float64) float64 {
	return (desc - bybit*1200) / desc * 100
}

func sendMonitorMessage(env *envs) {
	if env.Monitor.Token == "" || env.Monitor.ChatID == 0 {
		log.Println("Key is invalid for monitor")
		return
	}

	monitorBot, err := tgbotapi.NewBotAPI(env.Monitor.Token)
	if err != nil {
		panic(err)
	}

	monitorBot.Debug = false

	cnt := 0
	for {
		time.Sleep(env.Period["monitor"])
		totalPrices := readPrices()

		info := ""
		if cnt%5 == 0 {
			info += "http://home.5004.pe.kr:8080\n" +
				"KRWUSD:" + strconv.FormatFloat(totalPrices.Currency.Price[0].Price, 'f', -1, 64) +
				" FixKRWUSD: 1200\n\n"
			cnt = 0
		}

		info += "BTC: Bybit[" + strconv.FormatFloat(totalPrices.BybitPrice.Price[0].Price, 'f', -1, 64) + "]" +
			" Fund: " + strconv.FormatFloat(totalPrices.BybitPrice.Price[0].FundingRate, 'f', -1, 64) + "\n" +
			"   Upbit[" + strconv.FormatFloat(premiumRate(totalPrices.BybitPrice.Price[0].Price, totalPrices.UpbitPrice.Price[0].Price), 'f', 3, 64) + "%]" +
			" Bithumb[" + strconv.FormatFloat(premiumRate(totalPrices.BybitPrice.Price[0].Price, totalPrices.BithumbPrice.Price[0].Price), 'f', 3, 64) + "%]\n" +
			"ETH: Bybit[" + strconv.FormatFloat(totalPrices.BybitPrice.Price[1].Price, 'f', -1, 64) + "]" +
			" Fund: " + strconv.FormatFloat(totalPrices.BybitPrice.Price[1].FundingRate, 'f', -1, 64) + "\n" +
			"   Upbit[" + strconv.FormatFloat(premiumRate(totalPrices.BybitPrice.Price[1].Price, totalPrices.UpbitPrice.Price[1].Price), 'f', 3, 64) + "%]" +
			" Bithumb[" + strconv.FormatFloat(premiumRate(totalPrices.BybitPrice.Price[1].Price, totalPrices.BithumbPrice.Price[1].Price), 'f', 3, 64) + "%]"

		msg := tgbotapi.NewMessage(env.Monitor.ChatID, info)
		monitorBot.Send(msg)
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

	totalPrices.CreatedAt = time.Now().Unix()
	return
}

func allPrices(c *gin.Context) {
	totalPrices := readPrices()
	c.JSON(http.StatusOK, totalPrices)
}

func setEnvs(env *envs) {
	godotenv.Load()
	// Default Value
	env.Period["upbit"] = 4 * time.Second
	env.Period["bithumb"] = 5 * time.Second
	env.Period["bybit"] = 3 * time.Second
	env.Period["currency"] = 600 * time.Second
	env.Period["alarm"] = 10 * time.Second
	env.Period["monitor"] = 10 * time.Second

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
	alarmPeriod, _ := strconv.Atoi(os.Getenv("AlarmPeriodSeconds"))
	if alarmPeriod > 0 {
		env.Period["alarm"] = time.Duration(alarmPeriod) * time.Second
	}
	monitorPeriod, _ := strconv.Atoi(os.Getenv("MonitorPeriodSeconds"))
	if monitorPeriod > 0 {
		env.Period["monitor"] = time.Duration(monitorPeriod) * time.Second
	}

	env.Monitor.ChatID, _ = strconv.ParseInt(os.Getenv("MonitorChatID"), 10, 64)
	env.Monitor.Token = os.Getenv("MonitorToken")
	env.Alarm.ChatID, _ = strconv.ParseInt(os.Getenv("AlarmChatID"), 10, 64)
	env.Alarm.Token = os.Getenv("AlarmToken")
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
	router.Run()

	log.Fatal(http.ListenAndServe(":8080", router))
}
