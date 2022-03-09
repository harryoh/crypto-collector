package currency

import (
	"net/http"

	"github.com/harryoh/crypto-collector/exchange/currency/types"
	"github.com/harryoh/crypto-collector/util"
)

const (
	baseURL = "https://exchange.jaeheon.kr:23490"
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

// CurrencyRate :
func (client *Client) CurrencyRate(
	currency string,
) (rate *types.CurrencyRate, err error) {
	if !isValidCurrency(currency) {
		err = &InvalidParams{
			message: "Invalid currency",
		}
		return
	}

	options := &util.RequestOptions{URL: baseURL + "/query/" + currency}
	err = util.Request(client.httpClient, options, &rate)
	return
}
