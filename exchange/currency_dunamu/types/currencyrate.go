package types

// [
// 	{
// 		"code": "FRX.KRWUSD",
// 		"currencyCode": "USD",
// 		"currencyName": "달러",
// 		"country": "미국",
// 		"name": "미국 (KRW/USD)",
// 		"date": "2022-03-08",
// 		"time": "20:01:38",
// 		"recurrenceCount": 419,
// 		"basePrice": 1235.5,
// 		"openingPrice": 1225.4,
// 		"highPrice": 1238.6,
// 		"lowPrice": 1225.4,
// 		"change": "RISE",
// 		"changePrice": 5,
// 		"cashBuyingPrice": 1257.12,
// 		"cashSellingPrice": 1213.88,
// 		"ttBuyingPrice": 1223.4,
// 		"ttSellingPrice": 1247.6,
// 		"tcBuyingPrice": null,
// 		"fcSellingPrice": null,
// 		"exchangeCommission": 2.1037,
// 		"usDollarRate": 1,
// 		"high52wPrice": 1230.5,
// 		"high52wDate": "2022-03-07",
// 		"low52wPrice": 1105.2,
// 		"low52wDate": "2021-06-01",
// 		"currencyUnit": 1,
// 		"provider": "하나은행",
// 		"timestamp": 1646737299066,
// 		"id": 79,
// 		"createdAt": "2016-10-21T06:13:34.000+0000",
// 		"modifiedAt": "2022-03-08T12:09:28.000+0000",
// 		"signedChangePrice": 5,
// 		"signedChangeRate": 0.0040633889,
// 		"changeRate": 0.0040633889
// 	}
// ]

type CurrencyRate []struct {
	Update string  `json:"modifiedAt"`
	USDKRW float64 `json:"basePrice"`
}
