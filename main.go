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
	Rules        []rule
	CreatedAt    int64
}

type telegramKey struct {
	ChatID int64
	Token  string
}

type rule struct {
	Use      bool
	Symbol   string
	Exchange string
	AlarmMax float64
	AlarmMin float64
}

// // Rules :
// type Rules struct {
// 	rule []rule
// }

type envs struct {
	Period         map[string]time.Duration
	Monitor        telegramKey
	Alarm          telegramKey
	CurrencyAPIKey string
	Rules          []rule
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
	sendMsg := true
	var lastAlarmTimestamp int64 = 0

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

	for {
		var data *cache2go.CacheItem
		data, err = cache.Value("rule")
		if err != nil {
			fmt.Println(err)
			continue
		}
		env.Rules = data.Data().([]rule)

		totalPrices := readPrices()

		if len(totalPrices.BybitPrice.Price) < 1 {
			fmt.Println("Bybit Prices is NULL!")
			continue
		}

		bybitBTCStr := "[Bybit]" +
			" BTC:" + strconv.FormatFloat(totalPrices.BybitPrice.Price[0].Price, 'f', -1, 64)
		bybitETHStr := "[Bybit]" +
			" ETH:" + strconv.FormatFloat(totalPrices.BybitPrice.Price[1].Price, 'f', -1, 64)
		content := ""

		var premiumRateBithumbBTC float64 = 0
		var premiumRateBithumbETH float64 = 0
		var premiumRateUpbitBTC float64 = 0
		var premiumRateUpbitETH float64 = 0
		for _, r := range env.Rules {
			if r.Use != true {
				continue
			}
			// fmt.Println(r.Use, r.Symbol, r.Exchange, r.AlarmMin, r.AlarmMax, premiumRateUpbitBTC, premiumRateBithumbETH, premiumRateUpbitBTC, premiumRateUpbitETH)
			ruleText := "[RULE]: " + r.Exchange + r.Symbol +
				" Min:" + strconv.FormatFloat(r.AlarmMin, 'f', -1, 64) +
				" Max:" + strconv.FormatFloat(r.AlarmMax, 'f', -1, 64) + "\n"

			if (r.Exchange == "all" || r.Exchange == "bithumb") && len(totalPrices.BithumbPrice.Price) > 1 {
				if r.Symbol == "all" || r.Symbol == "BTC" {
					premiumRateBithumbBTC = premiumRate(totalPrices.BybitPrice.Price[0].Price, totalPrices.BithumbPrice.Price[0].Price)
					if premiumRateBithumbBTC <= r.AlarmMin || premiumRateBithumbBTC >= r.AlarmMax {
						if lastAlarmTimestamp+int64(env.Period["alarm"]/time.Second) > time.Now().Unix() {
							continue
						}
						content = bybitBTCStr
						content += "\n[Bithumb]" +
							" BTC:" + strconv.FormatFloat(totalPrices.BithumbPrice.Price[0].Price, 'f', -1, 64) +
							"(" + strconv.FormatFloat(premiumRateBithumbBTC, 'f', 3, 64) + "%)"
						msg := tgbotapi.NewMessage(env.Alarm.ChatID, ruleText+content)
						lastAlarmTimestamp = time.Now().Unix()
						if sendMsg == true {
							alarmBot.Send(msg)
						} else {
							fmt.Println(ruleText + content)
						}
					}
				}

				if r.Symbol == "all" || r.Symbol == "ETH" {
					premiumRateBithumbETH = premiumRate(totalPrices.BybitPrice.Price[1].Price, totalPrices.BithumbPrice.Price[1].Price)
					if premiumRateBithumbETH <= r.AlarmMin || premiumRateBithumbETH >= r.AlarmMax {
						if lastAlarmTimestamp+int64(env.Period["alarm"]/time.Second) > time.Now().Unix() {
							continue
						}
						content = bybitETHStr
						content += "\n[Bithumb]" +
							" ETH:" + strconv.FormatFloat(totalPrices.BithumbPrice.Price[1].Price, 'f', -1, 64) +
							"(" + strconv.FormatFloat(premiumRateBithumbETH, 'f', 3, 64) + "%)"
						msg := tgbotapi.NewMessage(env.Alarm.ChatID, ruleText+content)
						lastAlarmTimestamp = time.Now().Unix()
						if sendMsg == true {
							alarmBot.Send(msg)
						} else {
							fmt.Println(ruleText + content)
						}
					}
				}
			}

			if (r.Exchange == "all" || r.Exchange == "upbit") && len(totalPrices.UpbitPrice.Price) > 1 {
				if r.Symbol == "all" || r.Symbol == "BTC" {
					premiumRateUpbitBTC = premiumRate(totalPrices.BybitPrice.Price[0].Price, totalPrices.UpbitPrice.Price[0].Price)
					if premiumRateUpbitBTC <= r.AlarmMin || premiumRateUpbitBTC >= r.AlarmMax {
						if lastAlarmTimestamp+int64(env.Period["alarm"]/time.Second) > time.Now().Unix() {
							continue
						}
						content = bybitBTCStr
						content += "\n[Upbit]" +
							" BTC:" + strconv.FormatFloat(totalPrices.UpbitPrice.Price[0].Price, 'f', -1, 64) +
							"(" + strconv.FormatFloat(premiumRateUpbitBTC, 'f', 3, 64) + "%)"
						msg := tgbotapi.NewMessage(env.Alarm.ChatID, ruleText+content)
						lastAlarmTimestamp = time.Now().Unix()
						if sendMsg == true {
							alarmBot.Send(msg)
						} else {
							fmt.Println(ruleText + content)
						}
					}
				}

				if r.Symbol == "all" || r.Symbol == "ETH" {
					premiumRateUpbitETH = premiumRate(totalPrices.BybitPrice.Price[1].Price, totalPrices.UpbitPrice.Price[1].Price)
					if premiumRateUpbitETH <= r.AlarmMin || premiumRateUpbitETH >= r.AlarmMax {
						if lastAlarmTimestamp+int64(env.Period["alarm"]/time.Second) > time.Now().Unix() {
							continue
						}
						content = bybitETHStr
						content += "\n[Upbit]" +
							" ETH:" + strconv.FormatFloat(totalPrices.UpbitPrice.Price[1].Price, 'f', -1, 64) +
							"(" + strconv.FormatFloat(premiumRateUpbitETH, 'f', 3, 64) + "%)"
						msg := tgbotapi.NewMessage(env.Alarm.ChatID, ruleText+content)
						lastAlarmTimestamp = time.Now().Unix()
						if sendMsg == true {
							alarmBot.Send(msg)
						} else {
							fmt.Println(ruleText + content)
						}
					}
				}
			}
		}
		time.Sleep(time.Second * 1)
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
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
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

	c.JSON(200, gin.H{"code": 20000, "data": data.Data()})
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
		totalPrices.Rules = data.Data().([]rule)
	}

	totalPrices.CreatedAt = time.Now().Unix()
	return
}

func allPrices(c *gin.Context) {
	totalPrice := readPrices()
	c.JSON(http.StatusOK, gin.H{"code": 20000, "data": totalPrice})
}

func setRule(c *gin.Context) {
	var json []rule
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cache := _cache()
	cache.Add("rule", 0, json)

	c.JSON(http.StatusOK, gin.H{"code": 20000, "data": "OK"})
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

	ruleUse, _ := strconv.ParseBool(os.Getenv("RuleAlarmUse"))
	ruleAlarmMax, _ := strconv.ParseFloat(os.Getenv("RuleAlarmMax"), 64)
	ruleAlarmMin, _ := strconv.ParseFloat(os.Getenv("RuleAlarmMin"), 64)
	env.Rules = append(env.Rules, rule{
		Use:      ruleUse,
		Symbol:   "ETH",
		Exchange: "upbit",
		AlarmMax: ruleAlarmMax,
		AlarmMin: ruleAlarmMin,
	})
	env.Rules = append(env.Rules, rule{
		Use:      ruleUse,
		Symbol:   "ETH",
		Exchange: "bithumb",
		AlarmMax: ruleAlarmMax,
		AlarmMin: ruleAlarmMin,
	})
	cache := _cache()
	cache.Add("rule", 0, env.Rules)
}

// func initializeAppWithServiceAccount() *firebase.App {
// 	opt := option.WithCredentialsFile("./serviceAccountKey.json")
// 	app, err := firebase.NewApp(context.Background(), nil, opt)
// 	if err != nil {
// 		log.Fatalf("error initializing app: %v\n", err)
// 	}
// 	return app
// }

func main() {
	cache := _cache()

	env := &envs{
		Period: make(map[string]time.Duration),
	}

	// app := initializeAppWithServiceAccount()
	// client, err := app.Auth(context.Background())
	// if err != nil {
	// 	log.Fatalf("error getting Auth client: %v\n", err)
	// }
	// idToken := "eyJhbGciOiJSUzI1NiIsImtpZCI6IjYxMDgzMDRiYWRmNDc1MWIyMWUwNDQwNTQyMDZhNDFkOGZmMWNiYTgiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoi7Jik7ZqM6re8IiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hLS9BT2gxNEdoSGkzWG9pSEtQai1aVEM1Z3JEc29qUGJuOEt2dmJaNlZ1UWloTDhnUT1zOTYtYyIsImlzcyI6Imh0dHBzOi8vc2VjdXJldG9rZW4uZ29vZ2xlLmNvbS9jcnlwdG8tbWFuYWdlci1mZjc4YSIsImF1ZCI6ImNyeXB0by1tYW5hZ2VyLWZmNzhhIiwiYXV0aF90aW1lIjoxNjEzNDY3NzA2LCJ1c2VyX2lkIjoiMWc4cWFxNjc5bllxWE9lSmtZSFRLc21nMUw5MiIsInN1YiI6IjFnOHFhcTY3OW5ZcVhPZUprWUhUS3NtZzFMOTIiLCJpYXQiOjE2MTM0Njc3MDYsImV4cCI6MTYxMzQ3MTMwNiwiZW1haWwiOiJoYXJyeUA1MDA0LnBlLmtyIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImZpcmViYXNlIjp7ImlkZW50aXRpZXMiOnsiZ29vZ2xlLmNvbSI6WyIxMDk4NDUxMzg0ODA3Nzk1NzgyMDMiXSwiZW1haWwiOlsiaGFycnlANTAwNC5wZS5rciJdfSwic2lnbl9pbl9wcm92aWRlciI6Imdvb2dsZS5jb20ifX0.HZaMqiQTtegp4UPVRhenDrkBc6RfuyLmokQouHi4giVhdjqv6X89aRX7_udCAaoyeEn14-TLdpnposWQEZkht82t4ItAxQkdYdUI7Yn9lj9o9LKmfDes6IkddgvQ2iWbcwIs8bpPEUABs6Rr-fXCv15QJoDe7O6DK79ps8aLeRveROyC8QNszWco6bCFSlWevJUwt0ZSNsRO578asqnvFScyRv9p8D0bPw6blaqMJ9epcRWlcJvLdloTtxLvyrI8X7HdV3qVJroO15TrzuMsV0gKMxjADiuzc_H4bZGQa9MZ5pWbjr0maSshcXMcwYpRsYX1n1P7Oc-0JKQtkvuGaQ"
	// user, err := client.VerifyIDToken(context.Background(), idToken)
	// if err != nil {
	// 	log.Fatalf("error verifying ID token: %v\n", err)
	// }

	// log.Printf("Verified ID token: %v\n", user.Firebase.Identities["email"])

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
