package types

type ExchangeReq struct {
	Ip          string  `json:"ip,optional" form:"ip,optional"`
	Symbol      string  `json:"symbol,optional" form:"symbol,optional"`
	PageNo      int64   `json:"pageNo,optional" form:"pageNo,optional"`
	PageSize    int64   `json:"pageSize,optional" form:"pageSize,optional"`
	Price       float64 `json:"price,optional" form:"price,optional"`
	Amount      float64 `json:"amount,optional" form:"amount,optional"`
	Direction   string  `json:"direction,optional" form:"direction,optional"`
	Type        string  `json:"type,optional" form:"type,optional"`
	UseDiscount float64 `json:"useDiscount,optional" form:"useDiscount,optional"`
}

func (r *ExchangeReq) OrderValid() bool {
	if r.Direction == "" || r.Type == "" {
		return false
	}
	return true
}

type ExchangeOrder struct {
	Id            int64   `json:"id" from:"id"`
	OrderId       string  `json:"orderId" from:"orderId"`
	Amount        float64 `json:"amount" from:"amount"`
	BaseSymbol    string  `json:"baseSymbol" from:"baseSymbol"`
	CanceledTime  int64   `json:"canceledTime" from:"canceledTime"`
	CoinSymbol    string  `json:"coinSymbol" from:"coinSymbol"`
	CompletedTime int64   `json:"completedTime" from:"completedTime"`
	Direction     int     `json:"direction" from:"direction"`
	MemberId      int64   `json:"memberId" from:"memberId"`
	Price         string  `json:"price" from:"price"`
	Status        string  `json:"status" from:"status"`
	Symbol        string  `json:"symbol" from:"symbol"`
	Time          int64   `json:"time" from:"time"`
	TradedAmount  float64 `json:"tradedAmount" from:"tradedAmount"`
	Turnover      float64 `json:"turnover" from:"turnover"`
	Type          string  `json:"type" from:"type"`
	UseDiscount   string  `json:"useDiscount" from:"useDiscount"`
}
