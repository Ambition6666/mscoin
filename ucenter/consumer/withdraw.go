package consumer

import (
	"common/msdb"
	"encoding/json"
	"log"
	"ucenter/database"
	"ucenter/internal/domain"
	"ucenter/internal/model"
)

func WithdrawConsumer(kafkaCli *database.KafkaClient, db *msdb.MsDB) {
	withdrawDomain := domain.NewWithdrawDomain(db)
	for {
		kafkaData, err := kafkaCli.Read()
		if err != nil {
			log.Println(err)
			continue
		}
		var wr model.WithdrawRecord
		json.Unmarshal(kafkaData.Data, &wr)
		//调用btc rpc进行转账
		err = withdrawDomain.Withdraw(wr)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("提现成功")
		}
	}
}
