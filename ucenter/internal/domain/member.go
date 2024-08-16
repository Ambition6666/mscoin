package domain

import (
	"common/msdb"
	"common/tools"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"ucenter/internal/dao"
	"ucenter/internal/model"
	"ucenter/internal/repo"
)

type MemberDomain struct {
	repo.MemberRepo
}

func (d *MemberDomain) FindMemberByPhone(ctx context.Context, phone string) *model.Member {
	mem, err := d.MemberRepo.FindByPhone(ctx, phone)
	if err != nil {
		logx.Errorf("MemberDomain MemberDomain err: %s", err.Error())
		return nil
	}
	return mem
}

func (d *MemberDomain) FindMemberById(ctx context.Context, id int64) *model.Member {
	mem, err := d.MemberRepo.FindById(ctx, id)
	if err != nil {
		logx.Errorf("MemberDomain MemberDomain err: %s", err.Error())
		return nil
	}
	return mem
}

func (d *MemberDomain) Register(
	ctx context.Context,
	username string,
	phone string,
	password string,
	country string,
	promotion string,
	partner string) error {
	mem := model.NewMember()
	tools.Default(mem)
	mem.Id = 0

	//首先处理密码 密码要md5加密，但是md5不安全，所以我们给加上一个salt值

	salt, pwd := tools.Encode(password, nil)
	mem.Salt = salt
	mem.Password = pwd
	mem.MobilePhone = phone
	mem.Username = username
	mem.Country = country
	mem.PromotionCode = promotion
	mem.FillSuperPartner(partner)
	mem.MemberLevel = model.GENERAL
	mem.Avatar = "https://mszlu.oss-cn-beijing.aliyuncs.com/mscoin/defaultavatar.png"
	err := d.MemberRepo.Save(ctx, mem)
	return err
}

func (d *MemberDomain) UpdateLoginCount(id int64, incr int) {
	err := d.MemberRepo.UpdateLoginCount(context.Background(), id, incr)
	if err != nil {
		logx.Error(err)
	}
}

func NewMemberDomain(db *msdb.MsDB) *MemberDomain {
	return &MemberDomain{
		dao.NewMemberDao(db),
	}
}
