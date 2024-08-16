package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"grpc-common/ucenter/types/register"
	"time"
	"ucenter-api/internal/svc"
	types "ucenter-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	ctx, CancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer CancelFunc()

	regReq := &register.RegReq{}
	if err := copier.Copy(regReq, req); err != nil {
		return nil, err
	}
	_, err = l.svcCtx.URegisterRPC.RegisterByPhone(ctx, regReq)
	if err != nil {
		return nil, err
	}
	return
}

func (l *RegisterLogic) SendCode(req *types.CodeRequest) (resp *types.Response, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = l.svcCtx.URegisterRPC.SendCode(ctx, &register.CodeReq{
		Phone:   req.Phone,
		Country: req.Country,
	})
	if err != nil {
		return nil, err
	}
	return
}
