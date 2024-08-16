package types

type Request struct {
	Username     string      `json:"username"`
	Password     string      `json:"password,optional"`
	Captcha      *CaptchaReq `json:"captcha,optional"`
	Phone        string      `json:"phone,optional"`
	Promotion    string      `json:"promotion,optional"`
	Code         string      `json:"code,optional"`
	Country      string      `json:"country,optional"`
	SuperPartner string      `json:"superPartner,optional"`
	Ip           string      `json:"ip,optional"`
}

type CaptchaReq struct {
	Server string `json:"server"`
	Token  string `json:"token"`
}

type Response struct {
	Message string `json:"message"`
}

type CodeRequest struct {
	Phone   string `json:"phone,optional"`
	Country string `json:"country,optional"`
}

type NoRes struct {
}

type LoginReq struct {
	Username string      `json:"username"`
	Password string      `json:"password"`
	Captcha  *CaptchaReq `json:"captcha,optional"`
	Ip       string      `json:"ip,optional"`
}

type LoginRes struct {
	Username      string `json:"username"`
	Token         string `json:"token"`
	MemberLevel   string `json:"memberLevel"`
	RealName      string `json:"realName"`
	Country       string `json:"country"`
	Avatar        string `json:"avatar"`
	PromotionCode string `json:"promotionCode"`
	Id            int64  `json:"id"`
	LoginCount    int    `json:"loginCount"`
	SuperPartner  string `json:"superPartner"`
	MemberRate    int    `json:"memberRate"`
}

type AssetReq struct {
	CoinName  string `json:"coinName,optional" path:"coinName,optional"`
	Ip        string `json:"ip,optional"`
	Unit      string `json:"unit,optional" form:"unit,optional"`
	PageNo    int    `json:"pageNo,optional" form:"pageNo,optional"`
	PageSize  int    `json:"pageSize,optional" form:"pageSize,optional"`
	StartTime string `json:"startTime,optional" form:"startTime,optional"`
	EndTime   string `json:"endTime,optional" form:"endTime,optional"`
	Symbol    string `json:"symbol,optional" form:"symbol,optional"`
	Type      string `json:"type,optional" form:"type,optional"`
}

type MemberTransaction struct {
	Id          int64   `json:"id" from:"id"`
	Address     string  `json:"address" from:"address"`
	Amount      float64 `json:"amount" from:"amount"`
	CreateTime  string  `json:"createTime" from:"createTime"`
	Fee         float64 `json:"fee" from:"fee"`
	Flag        int     `json:"flag" from:"flag"`
	MemberId    int64   `json:"memberId" from:"memberId"`
	Symbol      string  `json:"symbol" from:"symbol"`
	Type        string  `json:"type" from:"type"`
	DiscountFee string  `json:"discountFee" from:"discountFee"`
	RealFee     string  `json:"realFee" from:"realFee"`
}

type Coin struct {
	Id                int     `json:"id" from:"id"`
	Name              string  `json:"name" from:"name"`
	CanAutoWithdraw   int     `json:"canAutoWithdraw" from:"canAutoWithdraw"`
	CanRecharge       int     `json:"canRecharge" from:"canRecharge"`
	CanTransfer       int     `json:"canTransfer" from:"canTransfer"`
	CanWithdraw       int     `json:"canWithdraw" from:"canWithdraw"`
	CnyRate           float64 `json:"cnyRate" from:"cnyRate"`
	EnableRpc         int     `json:"enableRpc" from:"enableRpc"`
	IsPlatformCoin    int     `json:"isPlatformCoin" from:"isPlatformCoin"`
	MaxTxFee          float64 `json:"maxTxFee" from:"maxTxFee"`
	MaxWithdrawAmount float64 `json:"maxWithdrawAmount" from:"maxWithdrawAmount"`
	MinTxFee          float64 `json:"minTxFee" from:"minTxFee"`
	MinWithdrawAmount float64 `json:"minWithdrawAmount" from:"minWithdrawAmount"`
	NameCn            string  `json:"nameCn" from:"nameCn"`
	Sort              int     `json:"sort" from:"sort"`
	Status            int     `json:"status" from:"status"`
	Unit              string  `json:"unit" from:"unit"`
	UsdRate           float64 `json:"usdRate" from:"usdRate"`
	WithdrawThreshold float64 `json:"withdrawThreshold" from:"withdrawThreshold"`
	HasLegal          int     `json:"hasLegal" from:"hasLegal"`
	ColdWalletAddress string  `json:"coldWalletAddress" from:"coldWalletAddress"`
	MinerFee          float64 `json:"minerFee" from:"minerFee"`
	WithdrawScale     int     `json:"withdrawScale" from:"withdrawScale"`
	AccountType       int     `json:"accountType" from:"accountType"`
	DepositAddress    string  `json:"depositAddress" from:"depositAddress"`
	Infolink          string  `json:"infolink" from:"infolink"`
	Information       string  `json:"information" from:"information"`
	MinRechargeAmount float64 `json:"minRechargeAmount" from:"minRechargeAmount"`
}

type MemberWallet struct {
	Id             int64   `json:"id" from:"id"`
	Address        string  `json:"address" from:"address"`
	Balance        float64 `json:"balance" from:"balance"`
	FrozenBalance  float64 `json:"frozenBalance" from:"frozenBalance"`
	ReleaseBalance float64 `json:"releaseBalance" from:"releaseBalance"`
	IsLock         int     `json:"isLock" from:"isLock"`
	MemberId       int64   `json:"memberId" from:"memberId"`
	Version        int     `json:"version" from:"version"`
	Coin           Coin    `json:"coin" from:"coinId"`
	ToReleased     float64 `json:"toReleased" from:"toReleased"`
}

type ApproveReq struct {
}
type MemberSecurity struct {
	Username             string `json:"username"`
	Id                   int64  `json:"id"`
	CreateTime           string `json:"createTime"`
	RealVerified         string `json:"realVerified"`  //是否实名认证
	EmailVerified        string `json:"emailVerified"` //是否有邮箱
	PhoneVerified        string `json:"phoneVerified"` //是否有手机号
	LoginVerified        string `json:"loginVerified"`
	FundsVerified        string `json:"fundsVerified"` //是否有交易密码
	RealAuditing         string `json:"realAuditing"`
	MobilePhone          string `json:"mobilePhone"`
	Email                string `json:"email"`
	RealName             string `json:"realName"`
	RealNameRejectReason string `json:"realNameRejectReason"`
	IdCard               string `json:"idCard"`
	Avatar               string `json:"avatar"`
	AccountVerified      string `json:"accountVerified"`
}
type WithdrawReq struct {
	Unit       string  `json:"unit,optional" form:"unit,optional"`
	Address    string  `json:"address,optional" form:"address,optional"`
	Amount     float64 `json:"amount,optional" form:"amount,optional"`
	Fee        float64 `json:"fee,optional" form:"fee,optional"`
	JyPassword string  `json:"jyPassword,optional" form:"jyPassword,optional"`
	Code       string  `json:"code,optional" form:"code,optional"`
	Page       int     `json:"page,optional" form:"page,optional"`
	PageSize   int     `json:"pageSize,optional" form:"pageSize,optional"`
}
type WithdrawWalletInfo struct {
	Unit            string          `json:"unit"`
	Threshold       float64         `json:"threshold"` //阈值
	MinAmount       float64         `json:"minAmount"` //最小提币数量
	MaxAmount       float64         `json:"maxAmount"` //最大提币数量
	MinTxFee        float64         `json:"minTxFee"`  //最小交易手续费
	MaxTxFee        float64         `json:"maxTxFee"`
	NameCn          string          `json:"nameCn"`
	Name            string          `json:"name"`
	Balance         float64         `json:"balance"`
	CanAutoWithdraw string          `json:"canAutoWithdraw"` //true false
	WithdrawScale   int             `json:"withdrawScale"`
	AccountType     int             `json:"accountType"`
	Addresses       []AddressSimple `json:"addresses"`
}

type AddressSimple struct {
	Remark  string `json:"remark"`
	Address string `json:"address"`
}