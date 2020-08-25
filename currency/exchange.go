package currency

import (
	"net/http"

	"github.com/harryoh/crypto-collector/currency/types"
	"github.com/harryoh/crypto-collector/util"
)

const (
	baseURL = "https://earthquake.kr:23490"
)

// InvalidParams :
type InvalidParams struct {
	message string
	Err     error
}

func (e *InvalidParams) Error() string {
	return e.message
}

// Client :
type Client struct {
	httpClient *http.Client
}

// NewClient :
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

// ExchangeRate :
func (client *Client) ExchangeRate(
	currency string,
) (usdkrw *types.RateUSDKRW, err error) {
	if !isValidCurrency(currency) {
		err = &InvalidParams{
			message: "Invalid currency",
		}
		return
	}

	options := &util.RequestOptions{URL: baseURL + "/query/" + currency}
	err = util.Request(client.httpClient, options, &usdkrw)
	return
}

// // LastPrice :
// func (client *Client) LastPrice() (candles []*types.MinuteCandle, err error) {
// 	return client.MinuteCandles(1, "KRW-BTC")
// }
