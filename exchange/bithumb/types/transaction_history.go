package types

// TransactionHistory :
type TransactionHistory struct {
	Status string `json:"status"`
	Data   []struct {
		TransactionDate string `json:"transaction_date"`
		Type            string `json:"type"`
		UnitsTraded     string `json:"units_traded"`
		Price           string `json:"price"`
		Total           string `json:"total"`
	} `json:"data"`
}
