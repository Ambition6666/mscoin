package member

import (
	"errors"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/net/context"
	"grpc-common/ucenter/types/member"
	"ucenter/internal/domain"
	"ucenter/internal/svc"
)

type MemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	MemberDomain *domain.MemberDomain
}

func NewMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberLogic {
	return &MemberLogic{
		ctx:          ctx,
		svcCtx:       svcCtx,
		Logger:       logx.WithContext(ctx),
		MemberDomain: domain.NewMemberDomain(svcCtx.Db),
	}
}

func (l *MemberLogic) FindMemberById(in *member.MemberReq) (*member.MemberRes, error) {
	ctx := context.Background()
	mem := l.MemberDomain.FindMemberById(ctx, in.MemberId)
	if mem == nil {
		return nil, errors.New("用户不存在")
	}
	resp := &member.MemberRes{}
	err := copier.Copy(resp, mem)

	if err != nil {
		logx.Errorf("通过id查询成员结果复制错误: %v", err)
		return nil, errors.New("查询错误")
	}

	return resp, nil
}
