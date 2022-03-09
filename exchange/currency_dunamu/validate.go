package currency

func isValidCurrency(currency string) bool {
	return currency == "USD_KRW" || currency == "KRW" || currency == "KRWUSD"
}
