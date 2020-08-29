package bithumb

import (
	"net/http"

	"github.com/harryoh/crypto-collector/exchange/bithumb/types"
	"github.com/harryoh/crypto-collector/util"
)

const (
	baseURL = "https://api.bithumb.com"
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

// TransactionHistory :
func (client *Client) TransactionHistory(
	symbol string,
) (txhistory *types.TransactionHistory, err error) {
	if !isValidSymbol(symbol) {
		err = &InvalidParams{
			message: "Invalid unit",
		}
		return
	}

	query := map[string]string{
		"count": "1",
	}

	options := &util.RequestOptions{
		URL:   baseURL + "/public/transaction_history/" + symbol,
		Query: query,
	}
	err = util.Request(client.httpClient, options, &txhistory)
	return
}

// LastPrice :
func (client *Client) LastPrice(
	symbol string,
) (txhistory *types.TransactionHistory, err error) {
	return client.TransactionHistory(symbol)
}
