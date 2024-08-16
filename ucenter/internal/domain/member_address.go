package domain

import (
	"common/msdb"
	"golang.org/x/net/context"
	"ucenter/internal/dao"
	"ucenter/internal/model"
	"ucenter/internal/repo"
)

type MemberAddressDomain struct {
	memberAddressRepo repo.MemberAddressRepo
}

func (d *MemberAddressDomain) FindAddressByCoinId(ctx context.Context, memId int64, coinId int64) ([]*model.MemberAddress, error) {
	return d.memberAddressRepo.FindByCoinId(ctx, memId, coinId)
}

func NewMemberAddressDomain(db *msdb.MsDB) *MemberAddressDomain {
	return &MemberAddressDomain{
		memberAddressRepo: dao.NewMemberAddressDao(db),
	}
}
