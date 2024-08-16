package domain

import (
	"common/btc"
	"common/msdb"
	"common/tools"
	"context"
	"errors"

	mclient "grpc-common/market/client"
	"grpc-common/market/types/market"
	"log"
	"ucenter/internal/dao"
	"ucenter/internal/model"
	"ucenter/internal/repo"
)

type WithdrawDomain struct {
	withdrawRecordRepo repo.WithdrawRecordRepo
	memberWalletDomain *MemberWalletDomain
	marketRPC          *mclient.Market
}

func NewWithdrawDomain(db *msdb.MsDB) *WithdrawDomain {
	return &WithdrawDomain{
		withdrawRecordRepo: dao.NewWithdrawDao(db),
		memberWalletDomain: NewMemberWalletDomain(db, nil),
	}
}
func (d *WithdrawDomain) SaveRecord(ctx context.Context, conn msdb.DbConn, wr *model.WithdrawRecord) error {
	return d.withdrawRecordRepo.Save(ctx, conn, wr)
}

func (d *WithdrawDomain) Withdraw(wr model.WithdrawRecord) error {
	//1. 获取到当前用户的address
	ctx := context.Background()
	wallet, err := d.memberWalletDomain.FindWalletByMemIdAndCoinId(ctx, wr.MemberId, wr.CoinId)
	if err != nil {
		return err
	}
	address := wallet.Address
	//2. 根据地址 获取到UTXO
	listUnspent, err := btc.ListUnspent(0, 999, []string{address})
	if err != nil {
		return err
	}
	//3. 判断UTXO中是否符合 提现金额 以及需要使用几个input
	totalAmount := wr.TotalAmount
	var inputs []btc.Input
	var utxoAmount float64
	for _, v := range listUnspent {
		utxoAmount += v.Amount
		input := btc.Input{Txid: v.Txid, Vout: v.Vout}
		inputs = append(inputs, input)
		if utxoAmount >= totalAmount {
			break
		}
	}
	if utxoAmount < totalAmount {
		return errors.New("账户余额不足")
	}
	//4. 创建交易
	//计算output  两个output 假设 utxo = 0.1 totalAmount = 0.002 arrivedAmount = 0.00185000 fee =0.00015000
	//fee = allinput - alloutput = utxoAmount（0.1） - 0.00185000 - (0.1-0.002)
	var outputs []map[string]any
	oneOutput := make(map[string]any)
	oneOutput[wr.Address] = wr.ArrivedAmount
	twoOutput := make(map[string]any)
	twoOutput[address] = tools.SubFloor(utxoAmount, totalAmount, 8)
	outputs = append(outputs, oneOutput, twoOutput)
	hexTxid, err := btc.CreateRawTransaction(inputs, outputs)
	if err != nil {
		return err
	}
	//5. 签名
	sign, err := btc.SignRawTransactionWithWallet(hexTxid)
	if err != nil {
		return err
	}
	//6. 广播交易
	txid, err := btc.SendRawTransaction(sign.Hex)
	if err != nil {
		return err
	}
	wr.TransactionNumber = txid
	wr.Status = 3 //提现成功
	err = d.withdrawRecordRepo.UpdateTransactionNumber(ctx, wr)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (d *WithdrawDomain) RecordList(
	ctx context.Context,
	userId int64,
	page int32,
	pageSize int32,
	marketRPC mclient.Market) ([]*model.WithdrawRecordVo, int64, error) {

	list, total, err := d.withdrawRecordRepo.FindByUserId(ctx, userId, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	var voList = make([]*model.WithdrawRecordVo, len(list))
	for i, v := range list {
		coin, err := marketRPC.FindCoinById(ctx, &market.MarketReq{
			Id: v.CoinId,
		})
		if err != nil {
			return nil, 0, err
		}
		voList[i] = v.ToVo(coin)
	}
	return voList, total, err
}
