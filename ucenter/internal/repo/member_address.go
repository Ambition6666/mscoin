package repo

import (
	"context"
	"ucenter/internal/model"
)

type MemberAddressRepo interface {
	Save(ctx context.Context, address *model.MemberAddress) error
	FindByCoinId(ctx context.Context, memberId int64, coinId int64) (list []*model.MemberAddress, err error)
	FindByMemberId(ctx context.Context, memberId int64) (list []*model.MemberAddress, err error)
}
