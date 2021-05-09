package util

// SymbolName :
func SymbolName(market *string) string {
	res := ""
	switch *market {
	case "BTC_KRW", "KRW-BTC", "BTCUSD":
		res = "BTC"
	case "ETH_KRW", "KRW-ETH", "ETHUSD":
		res = "ETH"
	case "XRP_KRW", "KRW-XRP", "XRPUSD":
		res = "XRP"
	case "EOS_KRW", "KRW-EOS", "EOSUSD":
		res = "EOS"
	// case "EOS_KRW", "KRW-EOS", "EOSUSD":
	// 	res = "EOS"
	case "USDKRW":
		res = "USDKRW"
	}

	return res
}
