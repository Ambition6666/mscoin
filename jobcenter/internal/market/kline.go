package market

import (
	"common/tools"
	"encoding/base64"
	"encoding/json"
	"jobcenter/internal/config"
	"jobcenter/internal/database"
	"jobcenter/internal/domain"
	"jobcenter/internal/model"
	"log"
	"strings"
	"sync"
	"time"
)

type Kline struct {
	okx         config.OkxConf
	wg          sync.WaitGroup
	klineDomain *domain.KlineDomain
	queueDomain *domain.QueueDomain
	redisClient *database.RedisClient
}

func NewKline(mcli *database.MongoClient, kcli *database.KafkaClient, conf config.OkxConf, rcli *database.RedisClient) *Kline {
	return &Kline{
		okx:         conf,
		wg:          sync.WaitGroup{},
		klineDomain: domain.NewKlineDomain(mcli),
		queueDomain: domain.NewQueueDomain(kcli),
		redisClient: rcli,
	}
}
func (k *Kline) Do(period string) {
	log.Println("============启动k线数据拉取==============")
	k.wg.Add(2)
	go k.syncToMongo("BTC-USDT", "BTC/USDT", period)
	go k.syncToMongo("ETH-USDT", "ETH/USDT", period)
	k.wg.Wait()
	log.Println("===============k线数据拉取结束===============")
}

func (k *Kline) syncToMongo(instId string, symbol, period string) {
	api := "GET/api/v5/market/candles?instId=" + instId + "&bar=" + period
	timestamp := tools.ISO(time.Now())
	sha256 := tools.ComputeHmacSha256(timestamp+api, k.okx.SecretKey)
	sign := base64.StdEncoding.EncodeToString([]byte(sha256))
	header := make(map[string]string)
	header["OK-ACCESS-KEY"] = k.okx.ApiKey
	header["OK-ACCESS-SIGN"] = sign
	header["OK-ACCESS-TIMESTAMP"] = timestamp
	header["OK-ACCESS-PASSPHRASE"] = k.okx.Pass
	respBody, err := tools.GetWithHeader(
		k.okx.Host+"/api/v5/market/candles?instId=BTC-USDT&bar="+period,
		header,
		k.okx.Proxy)
	if err != nil {
		log.Println(err)
	} else {

		resp := &model.OkxKlineRes{}
		err := json.Unmarshal(respBody, resp)
		if err != nil {
			log.Println(err)
		} else {
			if resp.Code == "0" {
				//代表成功
				k.klineDomain.Save(resp.Data, symbol, period)
				if "1m" == period {
					if len(resp.Data) > 0 {
						k.queueDomain.Sync1mKline(resp.Data[0], symbol, period)
					}
				}
			}
		}

		newString := strings.ReplaceAll(instId, "-", "::")

		k.redisClient.Cache.Set(newString+"::RATE", resp.Data[0][4])
	}
	k.wg.Done()

}

func (k *Kline) Send(data []string, symbol, period string) {
	if "1m" == period {
		//只有1m间隔的数据 才向kafka发送数据
		k.queueDomain.Sync1mKline(data, symbol, period)
	}
}
