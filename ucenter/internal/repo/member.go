package repo

import (
	"context"
	"ucenter/internal/model"
)

type MemberRepo interface {
	Save(ctx context.Context, mem *model.Member) error
	FindByPhone(ctx context.Context, phone string) (mem *model.Member, err error)
	FindById(ctx context.Context, id int64) (mem *model.Member, err error)
	UpdateLoginCount(ctx context.Context, id int64, incr int) error
}
