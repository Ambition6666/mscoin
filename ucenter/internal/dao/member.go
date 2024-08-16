package dao

import (
	"common/msdb"
	"common/msdb/gorms"
	"context"
	"errors"
	"fmt"
	"ucenter/internal/model"

	"gorm.io/gorm"
)

type MemberDao struct {
	conn *gorms.GormConn
}

func (m *MemberDao) Save(ctx context.Context, mem *model.Member) error {
	session := m.conn.Session(ctx)
	err := session.Save(mem).Error
	return err
}

func (m *MemberDao) FindById(ctx context.Context, id int64) (mem *model.Member, err error) {
	mem = model.NewMember()
	session := m.conn.Session(ctx)
	fmt.Println("连接", session)
	err = session.
		Model(&model.Member{}).
		Where("id = ?", id).
		Limit(1).
		Take(&mem).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func (m *MemberDao) FindByPhone(ctx context.Context, phone string) (mem *model.Member, err error) {
	mem = model.NewMember()
	session := m.conn.Session(ctx)
	fmt.Println("连接", session)
	err = session.
		Model(&model.Member{}).
		Where("mobile_phone = ?", phone).
		Limit(1).
		Take(&mem).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func (m *MemberDao) UpdateLoginCount(ctx context.Context, id int64, incr int) error {
	session := m.conn.Session(ctx)
	err := session.Exec("update member set login_count=login_count+? where id = ?", incr, id).Error
	return err
}

func NewMemberDao(db *msdb.MsDB) *MemberDao {
	return &MemberDao{
		conn: gorms.New(db.Conn),
	}
}
