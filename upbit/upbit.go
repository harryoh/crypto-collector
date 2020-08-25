package upbit

import (
	"net/http"
	"strconv"

	"github.com/harryoh/crypto-collector/upbit/types"
	"github.com/harryoh/crypto-collector/util"
)

const (
	baseURL = "https://api.upbit.com/v1"
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
	accessKey  string
	secretKey  string
	httpClient *http.Client
}

// NewClient :
func NewClient() *Client {
	return &Client{
		accessKey:  "_",
		secretKey:  "_",
		httpClient: &http.Client{},
	}
}

// Markets :
func (client *Client) Markets() (markets []*types.Market, err error) {
	options := &util.RequestOptions{URL: baseURL + "/market/all"}
	err = util.Request(client.httpClient, options, &markets)
	return
}

// MinuteCandles :
func (client *Client) MinuteCandles(
	unit int,
	market string,
	params ...map[string]string,
) (candles []*types.MinuteCandle, err error) {
	if !isValidMinuteCandleUnit(unit) {
		err = &InvalidParams{
			message: "Invalid unit",
		}
		return
	}

	query := map[string]string{
		"market": market,
	}

	if len(params) > 0 {
		for _, param := range params {
			for index, value := range param {
				query[index] = value
			}
		}
	}

	options := &util.RequestOptions{
		URL:   baseURL + "/candles/minutes/" + strconv.Itoa(unit),
		Query: query,
	}
	err = util.Request(client.httpClient, options, &candles)
	return
}

// LastPrice :
func (client *Client) LastPrice(
	market string,
) (candles []*types.MinuteCandle, err error) {
	return client.MinuteCandles(1, market)
}
