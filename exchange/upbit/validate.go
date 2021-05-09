package upbit

func isValidMinuteCandleUnit(unit int) bool {
	return unit == 1 || unit == 3 || unit == 5 || unit == 10 || unit == 15 ||
		unit == 30 || unit == 60 || unit == 240
}

func isValidSymbol(symbol string) bool {
	return symbol == "KRW-BTC" || symbol == "KRW-ETH" || symbol == "KRW-XRP" ||
		symbol == "KRW-EOS"
}
