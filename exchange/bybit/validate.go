package bybit

func isValidSymbol(symbol string) bool {
	return symbol == "BTCUSD" || symbol == "ETHUSD" || symbol == "EOSUSD" ||
		symbol == "XRPUSD"
}
