package registerlogic

import (
	"common/tools"
	"context"
	"errors"

	"grpc-common/ucenter/types/register"
	"time"
	"ucenter/internal/domain"
	"ucenter/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	CaptchaDomain *domain.CaptchaDomain
	MemberDomain  *domain.MemberDomain
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:           ctx,
		svcCtx:        svcCtx,
		Logger:        logx.WithContext(ctx),
		CaptchaDomain: domain.NewCaptchaDomain(),
		MemberDomain:  domain.NewMemberDomain(svcCtx.Db),
	}
}

func (l *RegisterLogic) RegisterByPhone(in *register.RegReq) (*register.RegRes, error) {
	// todo: add your logic here and delete this line
	isVerify := l.CaptchaDomain.Verify(
		l.svcCtx.Config.Captcha,
		in.Captcha.Server,
		in.Captcha.Token,
		2,
		in.Ip)
	if !isVerify {
		return nil, errors.New("人机验证不通过")
	}

	// 验证码
	//RedisCode := ""
	ctx := context.Background()
	//fmt.Println(l.svcCtx.Cache)
	//err := l.svcCtx.Cache.Get("RegisterRedisKey"+in.Phone, &RedisCode)
	//if err != nil {
	//	return nil, errors.New("验证码不可用或者验证码已过期")
	//}
	//if in.Code != RedisCode {
	//	return nil, errors.New("验证码不正确")
	//}
	//检查手机号是否注册
	mem := l.MemberDomain.FindMemberByPhone(ctx, in.Phone)
	if mem != nil {
		return nil, errors.New("手机号已经被注册")
	}
	err := l.MemberDomain.Register(
		ctx,
		in.Username,
		in.Phone,
		in.Password,
		in.Country,
		in.Promotion,
		in.SuperPartner,
	)
	if err != nil {
		return nil, errors.New("注册失败")
	}
	return &register.RegRes{}, nil
}

func (l *RegisterLogic) SendCode(in *register.CodeReq) (*register.NoRes, error) {
	code := tools.Rand4Num()
	l.Infof("验证码为: %s\n", code)
	l.Infof("手机号码为: %s\n, 国家区域为%s\n", in.Phone, in.Country)
	//通过短信平台发送验证码
	err := l.svcCtx.Cache.SetWithExpireCtx(context.Background(), "RegisterRedisKey"+in.Phone, code, 5*time.Minute)
	return &register.NoRes{}, err
}
