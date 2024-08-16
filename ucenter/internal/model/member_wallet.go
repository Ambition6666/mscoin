package model

import "grpc-common/market/types/market"

type MemberWallet struct {
	Id                int64   `gorm:"column:id"`
	Address           string  `gorm:"column:address"`
	AddressPrivateKey string  `gorm:"column:address_private_key"`
	Balance           float64 `gorm:"column:balance"`
	FrozenBalance     float64 `gorm:"column:frozen_balance"`
	ReleaseBalance    float64 `gorm:"column:release_balance"`
	IsLock            int     `gorm:"column:is_lock"`
	MemberId          int64   `gorm:"column:member_id"`
	Version           int     `gorm:"column:version"`
	CoinId            int64   `gorm:"column:coin_id"`
	ToReleased        float64 `gorm:"column:to_released"`
	CoinName          string  `gorm:"column:coin_name"`
}

func (*MemberWallet) TableName() string {
	return "member_wallet"
}

type MemberWalletCoin struct {
	Id             int64        `json:"id" from:"id"`
	Address        string       `json:"address" from:"address"`
	Balance        float64      `json:"balance" from:"balance"`
	FrozenBalance  float64      `json:"frozenBalance" from:"frozenBalance"`
	ReleaseBalance float64      `json:"releaseBalance" from:"releaseBalance"`
	IsLock         int          `json:"isLock" from:"isLock"`
	MemberId       int64        `json:"memberId" from:"memberId"`
	Version        int          `json:"version" from:"version"`
	Coin           *market.Coin `json:"coin" from:"coinId"`
	ToReleased     float64      `json:"toReleased" from:"toReleased"`
}

func NewMemberWallet(memId int64, coin *market.Coin) (*MemberWallet, *MemberWalletCoin) {
	mw := &MemberWallet{
		MemberId: memId,
		CoinId:   int64(coin.Id),
		CoinName: coin.Unit,
	}
	mwc := &MemberWalletCoin{
		MemberId: memId,
		Coin:     coin,
	}
	return mw, mwc
}
