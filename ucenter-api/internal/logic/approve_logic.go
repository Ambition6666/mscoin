package logic

import (
	"common/tools"
	"context"
	"grpc-common/ucenter/client"
	"time"
	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Approve struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Approve {
	return &Approve{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Approve) FindSecuritySetting(req *types.ApproveReq) (*types.MemberSecurity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	value := l.ctx.Value("userId").(int64)
	memberRes, err := l.svcCtx.UMemberRPC.FindMemberById(ctx, &client.MemberReq{
		MemberId: value,
	})
	if err != nil {
		return nil, err
	}
	ms := &types.MemberSecurity{}
	ms.Username = memberRes.Username
	ms.CreateTime = tools.ToTimeString(memberRes.RegistrationTime)
	ms.Id = memberRes.Id
	if memberRes.Email != "" {
		ms.EmailVerified = "true"
		ms.Email = memberRes.Email
	} else {
		ms.EmailVerified = "false"
	}
	if memberRes.JyPassword != "" {
		ms.FundsVerified = "true"
	} else {
		ms.FundsVerified = "false"
	}
	ms.LoginVerified = "true"
	if memberRes.MobilePhone != "" {
		ms.PhoneVerified = "true"
		ms.MobilePhone = memberRes.MobilePhone
	} else {
		ms.PhoneVerified = "false"
	}
	if memberRes.RealName != "" {
		ms.RealVerified = "true"
		ms.RealName = memberRes.RealName
	} else {
		ms.RealVerified = "false"
	}
	ms.IdCard = memberRes.IdNumber
	if memberRes.IdNumber != "" {
		ms.IdCard = memberRes.IdNumber[:2] + "********"
	}
	//0 未认证 1 审核中 2 已认证
	if memberRes.RealNameStatus == 1 {
		ms.RealAuditing = "true"
	} else {
		ms.RealAuditing = "false"
	}
	ms.Avatar = memberRes.Avatar
	if memberRes.Bank == "" && memberRes.AliNo == "" && memberRes.Wechat == "" {
		ms.AccountVerified = "false"
	} else {
		ms.AccountVerified = "true"
	}
	return ms, nil
}
