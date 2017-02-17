package client

import (
	"bytes"
	"crypto/md5"
	// "encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/milkbobo/gopay/common"
	"github.com/milkbobo/gopay/util"
	"sort"
	"strings"
	"time"
)

var defaultWechatAppClient *WechatAppClient

// DefaultWechatAppClient 默认微信app客户端
func DefaultWechatAppClient() *WechatAppClient {
	return defaultWechatAppClient
}

// WechatAppClient 微信app支付
type WechatAppClient struct {
	AppID       string // AppID
	MchID       string // 商户号ID
	CallbackURL string // 回调地址
	Key         string // 密钥
	PayURL      string // 支付地址
}

func InitWechatClient(c *WechatAppClient) {
	defaultWechatAppClient = c
}

// Pay 支付
func (wechat *WechatAppClient) Pay(charge *common.Charge) (map[string]string, error) {
	var m = make(map[string]string)
	m["appid"] = wechat.AppID
	m["mch_id"] = wechat.MchID
	m["nonce_str"] = util.RandomStr()
	m["body"] = charge.Describe
	m["out_trade_no"] = charge.TradeNum
	m["total_fee"] = fmt.Sprintf("%.2f", charge.MoneyFee)
	m["spbill_create_ip"] = util.LocalIP()
	m["notify_url"] = wechat.CallbackURL
	m["trade_type"] = "APP"

	sign, err := wechat.GenSign(m)
	if err != nil {
		return map[string]string{}, err
	}
	m["sign"] = sign
	// 转出xml结构
	buf := bytes.NewBufferString("")
	for k, v := range m {
		buf.WriteString(fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, v, k))
	}
	xmlStr := fmt.Sprintf("<xml>%s</xml>", buf.String())

	re, err := HTTPSC.PostData(wechat.PayURL, "text/xml:charset=UTF-8", xmlStr)
	if err != nil {
		return map[string]string{}, err
	}

	var xmlRe common.WeChatReResult
	err = xml.Unmarshal(re, &xmlRe)
	if err != nil {
		return map[string]string{}, err
	}

	if xmlRe.ReturnCode != "SUCCESS" {
		// 通信失败
		return map[string]string{}, errors.New(xmlRe.ReturnMsg)
	}

	if xmlRe.ResultCode != "SUCCESS" {
		// 支付失败
		return map[string]string{}, errors.New(xmlRe.ErrCodeDes)
	}

	var c = make(map[string]string)
	c["appid"] = wechat.AppID
	c["partnerid"] = wechat.MchID
	c["prepayid"] = xmlRe.PrepayID
	c["package"] = "Sign=WXPay"
	c["noncestr"] = util.RandomStr()
	c["timestamp"] = fmt.Sprintf("%d", time.Now().Unix())

	sign2, err := wechat.GenSign(c)
	if err != nil {
		return map[string]string{}, err
	}
	//c["signType"] = "MD5"
	c["paySign"] = strings.ToUpper(sign2)

	return c, nil
}

// GenSign 产生签名
func (wechat *WechatAppClient) GenSign(m map[string]string) (string, error) {
	delete(m, "sign")
	delete(m, "key")
	var signData []string
	for k, v := range m {
		if v != "" {
			signData = append(signData, fmt.Sprintf("%s=%s", k, v))
		}
	}

	sort.Strings(signData)
	signStr := strings.Join(signData, "&")
	signStr = signStr + "&key=" + wechat.Key
	c := md5.New()
	_, err := c.Write([]byte(signStr))
	if err != nil {
		return "", err
	}
	signByte := c.Sum(nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", signByte), nil
}

// CheckSign 检查签名
func (wechat *WechatAppClient) CheckSign(data string, sign string) error {
	signData := data + "&key=" + wechat.Key
	c := md5.New()
	_, err := c.Write([]byte(signData))
	if err != nil {
		return err
	}
	signOut := fmt.Sprintf("%x", c.Sum(nil))
	if strings.ToUpper(sign) == strings.ToUpper(signOut) {
		return nil
	}
	return errors.New("签名交易错误")
}
