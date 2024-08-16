package server

import (
	"context"

	"grpc-common/ucenter/types/member"
	memberlogic "ucenter/internal/logic/member"
	"ucenter/internal/svc"
)

type MemberServer struct {
	svcCtx *svc.ServiceContext
	member.UnimplementedMemberServer
}

func NewMemberServer(svcCtx *svc.ServiceContext) *MemberServer {
	return &MemberServer{
		svcCtx: svcCtx,
	}
}

func (s *MemberServer) FindMemberById(ctx context.Context, in *member.MemberReq) (*member.MemberRes, error) {
	l := memberlogic.NewMemberLogic(ctx, s.svcCtx)
	return l.FindMemberById(in)
}
