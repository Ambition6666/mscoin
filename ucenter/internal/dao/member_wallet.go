package dao

import (
	"common/msdb"
	"common/msdb/gorms"
	"context"
	"errors"
	"gorm.io/gorm"
	"ucenter/internal/model"
)

type MemberWalletDao struct {
	conn *gorms.GormConn
}

func (m *MemberWalletDao) Save(ctx context.Context, mw *model.MemberWallet) error {
	session := m.conn.Session(ctx)
	err := session.Save(&mw).Error
	return err
}

func (m *MemberWalletDao) FindByIdAndCoinName(ctx context.Context, memId int64, coinName string) (mw *model.MemberWallet, err error) {
	session := m.conn.Session(ctx)
	err = session.Model(&model.MemberWallet{}).
		Where("member_id=? and coin_name=?", memId, coinName).
		Take(&mw).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func (m *MemberWalletDao) FindByIdAndCoinId(ctx context.Context, memId int64, coinId int64) (mw *model.MemberWallet, err error) {
	session := m.conn.Session(ctx)
	err = session.Model(&model.MemberWallet{}).
		Where("member_id=? and coin_id=?", memId, coinId).
		Take(&mw).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func (m *MemberWalletDao) UpdateFreeze(ctx context.Context, conn msdb.DbConn, memId int64, money float64, symbol string) error {
	m.conn = conn.(*gorms.GormConn)
	session := m.conn.Tx(ctx)
	query := "update member_wallet set balance=balance-?,frozen_balance=frozen_balance+? where member_id=? and coin_name=? and balance > ?"
	exec := session.Exec(query, money, money, memId, symbol, money)
	err := exec.Error
	if err != nil {
		return err
	}
	affected := exec.RowsAffected
	if affected <= 0 {
		return errors.New("no update row")
	}
	return nil
}

func (m *MemberWalletDao) UpdateWallet(ctx context.Context, conn msdb.DbConn, id int64, balance float64, frozenBalance float64) error {
	gormConn := conn.(*gorms.GormConn)
	tx := gormConn.Tx(ctx)
	//Update
	updateSql := "update member_wallet set balance=?,frozen_balance=? where id=?"
	err := tx.Model(&model.MemberWallet{}).Exec(updateSql, balance, frozenBalance, id).Error
	return err
}

func (m *MemberWalletDao) FindByMemId(ctx context.Context, memId int64) (list []*model.MemberWallet, err error) {
	session := m.conn.Session(ctx)
	err = session.Model(&model.MemberWallet{}).Where("member_id=?", memId).Find(&list).Error
	return
}

func (m *MemberWalletDao) UpdateAddress(ctx context.Context, mw *model.MemberWallet) error {
	session := m.conn.Session(ctx)
	return session.Model(&model.MemberWallet{}).Where("id = ?", mw.Id).Update("address", mw.Address).Update("address_private_key", mw.AddressPrivateKey).Error
}
func (m *MemberWalletDao) GetAddress(ctx context.Context, name string) (list []string, err error) {
	session := m.conn.Session(ctx)
	err = session.Model(&model.MemberWallet{}).Where("coin_name = ?", name).Select("address").Find(&list).Error
	return
}
func NewMemberWalletDao(db *msdb.MsDB) *MemberWalletDao {
	return &MemberWalletDao{
		conn: gorms.New(db.Conn),
	}
}
