package server

import (
	"context"

	"grpc-common/ucenter/types/login"
	loginlogic "ucenter/internal/logic/login"
	"ucenter/internal/svc"
)

type LoginServer struct {
	svcCtx *svc.ServiceContext
	login.UnimplementedLoginServer
}

func NewLoginServer(svcCtx *svc.ServiceContext) *LoginServer {
	return &LoginServer{
		svcCtx: svcCtx,
	}
}

func (s *LoginServer) Login(ctx context.Context, in *login.LoginReq) (*login.LoginRes, error) {
	l := loginlogic.NewLoginLogic(ctx, s.svcCtx)
	return l.Login(in)
}
