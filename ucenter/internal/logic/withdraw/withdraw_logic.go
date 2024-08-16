package withdraw

import (
	"common/msdb"
	"common/msdb/tran"
	"common/tools"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/net/context"
	"grpc-common/ucenter/types/withdraw"
	"time"
	"ucenter/database"
	"ucenter/internal/domain"
	"ucenter/internal/model"
	"ucenter/internal/svc"
)

type WithdrawLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	memberAddressDomain *domain.MemberAddressDomain
	memberDomain        *domain.MemberDomain
	memberWalletDomain  *domain.MemberWalletDomain
	transaction         tran.Transaction
	withdrawDomain      *domain.WithdrawDomain
	kafkaCli            *database.KafkaClient
}

func (l *WithdrawLogic) FindAddressByCoinId(req *withdraw.WithdrawReq) (*withdraw.AddressSimpleList, error) {
	ctx := context.Background()
	userId := req.UserId
	coinId := req.CoinId
	memberAddresses, err := l.memberAddressDomain.FindAddressByCoinId(ctx, userId, coinId)
	if err != nil {
		return nil, err
	}
	var list []*withdraw.AddressSimple
	copier.Copy(&list, memberAddresses)
	return &withdraw.AddressSimpleList{
		List: list,
	}, nil
}

const WithdrawRedisKey = "WITHDRAW::"

func (l *WithdrawLogic) SendCode(in *withdraw.WithdrawReq) (*withdraw.NoRes, error) {
	code := "1234"
	logx.Infof("验证码为: %s", code)
	//通过短信平台发送验证码
	err := l.svcCtx.Cache.SetWithExpireCtx(context.Background(), WithdrawRedisKey+in.Phone, code, 5*time.Minute)
	return &withdraw.NoRes{}, err
}

func (l *WithdrawLogic) WithdrawCode(req *withdraw.WithdrawReq) (*withdraw.NoRes, error) {
	//1. 先判断code是否合法
	member := l.memberDomain.FindMemberById(l.ctx, req.UserId)
	var redisCode string
	err := l.svcCtx.Cache.GetCtx(l.ctx, WithdrawRedisKey+member.MobilePhone, &redisCode)
	if err != nil {
		return nil, err
	}
	if req.Code != redisCode {
		if err != nil {
			return nil, errors.New("验证码不正确")
		}
	}
	//2. 检查交易密码是否正确'
	if req.JyPassword != member.JyPassword {
		return nil, errors.New("交易密码错误")
	}
	//3. 检查钱包余额是否足够
	wallet, err := l.memberWalletDomain.FindWalletByMemIdAndCoinName(l.ctx, req.UserId, req.Unit)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, errors.New("钱包不存在")
	}
	if wallet.Balance < req.Amount {
		return nil, errors.New("余额不足")
	}
	//4.冻结
	err = l.transaction.Action(func(conn msdb.DbConn) error {
		err := l.memberWalletDomain.Freeze(l.ctx, conn, req.UserId, req.Amount, req.Unit)
		if err != nil {
			return err
		}
		//5.提现记录 状态设为正在审核中...
		wr := &model.WithdrawRecord{}
		wr.CoinId = req.CoinId
		wr.Address = req.Address
		wr.Fee = req.Fee
		wr.TotalAmount = req.Amount
		wr.ArrivedAmount = tools.SubFloor(req.Amount, req.Fee, 10)
		wr.Remark = ""
		wr.CanAutoWithdraw = 0
		wr.IsAuto = 0
		wr.Status = 0 //审核中
		wr.CreateTime = time.Now().UnixMilli()
		wr.DealTime = 0
		wr.MemberId = req.UserId
		wr.TransactionNumber = "" //目前还没有交易编号
		err = l.withdrawDomain.SaveRecord(l.ctx, conn, wr)
		if err != nil {
			return err
		}
		marshal, _ := json.Marshal(wr)
		//6. 发送MQ 处理提现事宜
		kafkaData := database.KafkaData{
			Topic: "withdraw",
			Key:   []byte(fmt.Sprintf("%d", req.UserId)),
			Data:  marshal,
		}
		//记得 wr.AllowAutoTopicCreation = true
		err = l.kafkaCli.SendSync(kafkaData)
		if err != nil {
			logx.Error(err)
			return errors.New("MQ发送失败")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &withdraw.NoRes{}, nil
}

func (l *WithdrawLogic) WithdrawRecord(req *withdraw.WithdrawReq) (*withdraw.RecordList, error) {
	list, total, err := l.withdrawDomain.RecordList(l.ctx, req.UserId, req.Page, req.PageSize, l.svcCtx.MarketRPC)
	if err != nil {
		return nil, err
	}
	var rList []*withdraw.WithdrawRecord
	copier.Copy(&rList, list)
	return &withdraw.RecordList{
		List:  rList,
		Total: total,
	}, nil
}

func NewWithdrawLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WithdrawLogic {
	return &WithdrawLogic{
		ctx:                 ctx,
		svcCtx:              svcCtx,
		Logger:              logx.WithContext(ctx),
		memberDomain:        domain.NewMemberDomain(svcCtx.Db),
		memberAddressDomain: domain.NewMemberAddressDomain(svcCtx.Db),
		memberWalletDomain:  domain.NewMemberWalletDomain(svcCtx.Db, svcCtx.Cache),
		withdrawDomain:      domain.NewWithdrawDomain(svcCtx.Db),
		transaction:         tran.NewTransaction(svcCtx.Db.Conn),
		kafkaCli:            svcCtx.Kcli,
	}
}
