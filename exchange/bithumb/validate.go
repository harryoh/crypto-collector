package bithumb

func isValidSymbol(symbol string) bool {
	return symbol == "BTC_KRW" || symbol == "ETH_KRW" || symbol == "XRP_KRW"
}
