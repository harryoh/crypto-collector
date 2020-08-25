package bybit

import (
	"net/http"

	"github.com/harryoh/crypto-collector/bybit/types"
	"github.com/harryoh/crypto-collector/util"
)

const (
	baseURL = "https://api.bybit.com/v2"
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

// Tickers :
func (client *Client) Tickers(
	symbol string,
) (tickers *types.Tickers, err error) {
	if !isValidSymbol(symbol) {
		err = &InvalidParams{
			message: "Invalid unit",
		}
		return
	}

	query := map[string]string{
		"symbol": symbol,
	}

	options := &util.RequestOptions{
		URL:   baseURL + "/public/tickers",
		Query: query,
	}
	err = util.Request(client.httpClient, options, &tickers)
	return
}

// LastPrice :
func (client *Client) LastPrice(
	symbol string,
) (tickers *types.Tickers, err error) {
	return client.Tickers(symbol)
}
