package dao

import (
	"common/msdb"
	"common/msdb/gorms"
	"golang.org/x/net/context"
	"ucenter/internal/model"
)

type MemberAddressDao struct {
	conn *gorms.GormConn
}

func (m *MemberAddressDao) Save(ctx context.Context, address *model.MemberAddress) error {
	session := m.conn.Session(ctx)
	err := session.Create(&address).Error
	return err
}

func (m *MemberAddressDao) FindByCoinId(ctx context.Context, memberId int64, coinId int64) (list []*model.MemberAddress, err error) {
	session := m.conn.Session(ctx)
	err = session.Model(&model.MemberAddress{}).
		Where("member_id=? and coin_id=? and status = ?", memberId, coinId, 0).
		Find(&list).Error
	return
}

func (m *MemberAddressDao) FindByMemberId(ctx context.Context, memberId int64) (list []*model.MemberAddress, err error) {
	session := m.conn.Session(ctx)
	err = session.Model(&model.MemberAddress{}).
		Where("member_id=? and status = ?", memberId, 0).
		Find(&list).Error
	return
}

func NewMemberAddressDao(db *msdb.MsDB) *MemberAddressDao {
	return &MemberAddressDao{
		conn: gorms.New(db.Conn),
	}
}
