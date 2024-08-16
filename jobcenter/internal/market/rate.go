package market

import (
	"common/tools"
	"encoding/base64"
	"encoding/json"
	"jobcenter/internal/database"
	"jobcenter/internal/model"
	"jobcenter/internal/svc"

	"log"
	"sync"
	"time"
)

type Rate struct {
	wg          sync.WaitGroup
	redisClient *database.RedisClient
}

// Do 获取人民币对美元汇率
func (r *Rate) Do() {
	r.wg.Add(1)
	go r.rateToRedis()
	r.wg.Wait()

}

func (r *Rate) rateToRedis() {
	api := "GET/api/v5/market/exchange-rate"
	timestamp := tools.ISO(time.Now())
	sha256 := tools.ComputeHmacSha256(timestamp+api, "secretKey")
	sign := base64.StdEncoding.EncodeToString([]byte(sha256))
	header := make(map[string]string)
	header["OK-ACCESS-KEY"] = "d5a748c6-214d-4fae-bef3-d32368ecbbe8"
	header["OK-ACCESS-SIGN"] = sign
	header["OK-ACCESS-TIMESTAMP"] = timestamp
	header["OK-ACCESS-PASSPHRASE"] = "Mszlu!@#$56789"
	respBody, err := tools.GetWithHeader(
		"https://www.okx.com/api/v5/market/exchange-rate",
		header,
		"http://127.0.0.1:7890")
	if err != nil {
		log.Println(err)
		r.wg.Done()
		return
	}
	//{
	//    "code": "0",
	//    "msg": "",
	//    "data": [ {
	//            "usdCny": "6.44"
	//}]
	//}
	resp := &model.OkxRateRes{}
	err = json.Unmarshal(respBody, resp)
	if err != nil {
		log.Println(err)
		r.wg.Done()
		return
	}
	for _, v := range resp.Data {
		r.redisClient.Cache.Set("USDT::CNY::RATE", v.UsdCny)
	}
	r.wg.Done()
}

func NewRate(sc *svc.ServiceContext) *Rate {
	return &Rate{
		redisClient: sc.RedisClient,
	}
}
