package parsers

type ParsContent struct {
	Source  string `json:"source"`
	Subject string `json:"subject"`
	Img     string `json:"img"`
	Text    string `json:"text"`
}

type CryptoCompareResponse struct {
	Data []struct {
		CoinInfo struct {
			Name     string `json:"Name"`
			FullName string `json:"FullName"`
		} `json:"CoinInfo"`
		Raw struct {
			Usd struct {
				Price           float64 `json:"PRICE"`
				ChangePCT24Hour float64 `json:"CHANGEPCT24HOUR"`
			} `json:"USD"`
		} `json:"RAW"`
	}
}
