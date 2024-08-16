package domain

import (
	"common/tools"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"ucenter/internal/config"
)

type CaptchaDomain struct {
}

type vaptchaReq struct {
	Id        string `json:"id"`
	Secretkey string `json:"secretkey"`
	Scene     int    `json:"scene"`
	Token     string `json:"token"`
	Ip        string `json:"ip"`
}
type vaptchaRsp struct {
	Success int    `json:"success"`
	Score   int    `json:"score"`
	Msg     string `json:"msg"`
}

func (d *CaptchaDomain) Verify(
	c config.CaptchaConf,
	server string,
	token string,
	scene int,
	ip string) bool {
	req := &vaptchaReq{
		Id:        c.Vid,
		Secretkey: c.Key,
		Scene:     scene,
		Token:     token,
		Ip:        ip,
	}
	respBytes, err := tools.Post(server, req)
	if err != nil {
		logx.Errorf("CaptchaDomain Verify post err : %s\n", err.Error())
		// fmt.Println(err)
		return false
	}
	var vaptchaRsp *vaptchaRsp
	err = json.Unmarshal(respBytes, &vaptchaRsp)
	if err != nil {
		logx.Errorf("CaptchaDomain Verify Unmarshal respBytes err : %s\n", err.Error())
		// fmt.Println(err)
		return false
	}

	if vaptchaRsp != nil && vaptchaRsp.Success == 1 {
		logx.Info("CaptchaDomain Verify no success\n")
		return true
	}
	return false
}

func NewCaptchaDomain() *CaptchaDomain {
	return &CaptchaDomain{}
}
