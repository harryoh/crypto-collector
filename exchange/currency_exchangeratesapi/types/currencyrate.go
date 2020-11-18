package types

// CurrencyRate :
type CurrencyRate struct {
	Rates struct {
		USDKRW float64 `json:"KRW"`
	} `json:"rates"`
	Base   string `json:"base"`
	Update string `json:"date"`
}
