package model

type Coin struct {
	Id                int64   `db:"id"`
	Name              string  `db:"name"`                // 货币
	CanAutoWithdraw   int64   `db:"can_auto_withdraw"`   // 是否能自动提币
	CanRecharge       int64   `db:"can_recharge"`        // 是否能充币
	CanTransfer       int64   `db:"can_transfer"`        // 是否能转账
	CanWithdraw       int64   `db:"can_withdraw"`        // 是否能提币
	CnyRate           float64 `db:"cny_rate"`            // 对人民币汇率
	EnableRpc         int64   `db:"enable_rpc"`          // 是否支持rpc接口
	IsPlatformCoin    int64   `db:"is_platform_coin"`    // 是否是平台币
	MaxTxFee          float64 `db:"max_tx_fee"`          // 最大提币手续费
	MaxWithdrawAmount float64 `db:"max_withdraw_amount"` // 最大提币数量
	MinTxFee          float64 `db:"min_tx_fee"`          // 最小提币手续费
	MinWithdrawAmount float64 `db:"min_withdraw_amount"` // 最小提币数量
	NameCn            string  `db:"name_cn"`             // 中文名称
	Sort              int64   `db:"sort"`                // 排序
	Status            int64   `db:"status"`              // 状态 0 正常 1非法
	Unit              string  `db:"unit"`                // 单位
	UsdRate           float64 `db:"usd_rate"`            // 对美元汇率
	WithdrawThreshold float64 `db:"withdraw_threshold"`  // 提现阈值
	HasLegal          int64   `db:"has_legal"`           // 是否是合法币种
	ColdWalletAddress string  `db:"cold_wallet_address"` // 冷钱包地址
	MinerFee          float64 `db:"miner_fee"`           // 转账时付给矿工的手续费
	WithdrawScale     int64   `db:"withdraw_scale"`      // 提币精度
	AccountType       int64   `db:"account_type"`        // 币种账户类型0：默认  1：EOS类型
	DepositAddress    string  `db:"deposit_address"`     // 充值地址
	Infolink          string  `db:"infolink"`            // 币种资料链接
	Information       string  `db:"information"`         // 币种简介
	MinRechargeAmount float64 `db:"min_recharge_amount"` // 最小充值数量
}
