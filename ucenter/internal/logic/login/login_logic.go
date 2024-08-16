package loginlogic

import (
	"common/tools"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/ucenter/types/login"
	"time"
	"ucenter/internal/domain"
	"ucenter/internal/svc"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	MemberDomain  *domain.MemberDomain
	CaptchaDomain *domain.CaptchaDomain
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:           ctx,
		svcCtx:        svcCtx,
		Logger:        logx.WithContext(ctx),
		MemberDomain:  domain.NewMemberDomain(svcCtx.Db),
		CaptchaDomain: domain.NewCaptchaDomain(),
	}
}

func (l *LoginLogic) Login(in *login.LoginReq) (*login.LoginRes, error) {
	//校验人机
	isVerify := l.CaptchaDomain.Verify(
		l.svcCtx.Config.Captcha,
		in.Captcha.Server,
		in.Captcha.Token,
		2,
		in.Ip)
	if !isVerify {
		return nil, errors.New("人机验证不通过")
	}
	//查询salt
	ctx := context.Background()
	mem := l.MemberDomain.FindMemberByPhone(ctx, in.Username)
	if mem == nil {
		return nil, errors.New("用户不存在")
	}
	salt := mem.Salt
	verify := tools.Verify(in.Password, salt, mem.Password, nil)
	if !verify {
		return nil, errors.New("账号密码不正确")
	}
	accessExpire := l.svcCtx.Config.JWT.AccessExpire
	accessSecret := l.svcCtx.Config.JWT.AccessSecret
	token, err := l.getJwtToken(accessSecret, time.Now().Unix(), accessExpire, mem.Id)
	if err != nil {
		return nil, errors.New("未知错误，请联系管理员")
	}
	loginCount := mem.LoginCount + 1
	go func() {
		l.MemberDomain.UpdateLoginCount(mem.Id, 1)
	}()
	return &login.LoginRes{
		Token:         token,
		Id:            mem.Id,
		Username:      mem.Username,
		MemberLevel:   mem.MemberLevelStr(),
		MemberRate:    mem.MemberRate(),
		RealName:      mem.RealName,
		Country:       mem.Country,
		Avatar:        mem.Avatar,
		PromotionCode: mem.PromotionCode,
		SuperPartner:  mem.SuperPartner,
		LoginCount:    int32(loginCount),
	}, nil
}

func (l *LoginLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
