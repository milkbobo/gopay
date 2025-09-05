package client

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/milkbobo/gopay/common"
	"net/url"
	"sort"
	"strings"
)

var aliWebClient *AliWebClient

// AliWebClient 支付宝网页支付
type AliWebClient struct {
	PartnerID   string          // 支付宝合作身份ID
	SellerID    string          // 卖家支付宝用户号
	AppID       string          // 支付宝分配给开发者的应用ID ps: 查询订单用
	CallbackURL string          // 回调接口
	PrivateKey  *rsa.PrivateKey // 私钥
	PublicKey   *rsa.PublicKey  // 公钥
}

func InitAliWebClient(c *AliWebClient) {
	aliWebClient = c
}

// DefaultAliWebClient 默认支付宝网页支付客户端
func DefaultAliWebClient() *AliWebClient {
	return aliWebClient
}

// Pay 实现支付接口
func (this *AliWebClient) Pay(charge *common.Charge) (map[string]string, error) {
	var m = make(map[string]string)
	m["service"] = "create_direct_pay_by_user"
	m["partner"] = this.PartnerID
	m["_input_charset"] = "UTF-8"
	m["notify_url"] = charge.CallbackURL
	m["return_url"] = charge.ReturnURL // 注意链接不能有&符号，否则会签名错误
	m["out_trade_no"] = charge.TradeNum
	m["subject"] = TruncatedText(charge.Describe, 32)
	m["total_fee"] = AliyunMoneyFeeToString(charge.MoneyFee)
	m["seller_id"] = this.SellerID

	sign := this.GenSign(m)

	m["sign"] = sign
	m["sign_type"] = "RSA"
	return map[string]string{"url": ToURL("https://mapi.alipay.com/gateway.do", m)}, nil
}

func (this *AliWebClient) CloseOrder(charge *common.Charge) (map[string]string, error) {
	return map[string]string{}, errors.New("暂未开发该功能")
}

func (this *AliWebClient) PayToClient(charge *common.Charge) (map[string]string, error) {
	return map[string]string{}, errors.New("暂未开发该功能")
}

// 订单查询
func (this *AliWebClient) QueryOrder(outTradeNo string) (common.AliWebQueryResult, error) {
	var m = make(map[string]string)
	m["service"] = "single_trade_query"
	m["partner"] = this.PartnerID
	m["_input_charset"] = "utf-8"
	m["out_trade_no"] = outTradeNo

	sign := this.GenSign(m)

	m["sign"] = sign
	m["sign_type"] = "RSA"
	return GetAlipay(ToURL("https://mapi.alipay.com/gateway.do", m))
}

// GenSign 产生签名
func (this *AliWebClient) GenSign(m map[string]string) string {
	var data []string
	for k, v := range m {
		if v != "" && k != "sign" && k != "sign_type" {
			data = append(data, fmt.Sprintf(`%s=%s`, k, v))
		}
	}
	sort.Strings(data)
	signData := strings.Join(data, "&")
	s := sha1.New()
	_, err := s.Write([]byte(signData))
	if err != nil {
		panic(err)
	}
	hashByte := s.Sum(nil)
	signByte, err := this.PrivateKey.Sign(rand.Reader, hashByte, crypto.SHA1)
	if err != nil {
		panic(err)
	}
	return url.QueryEscape(base64.StdEncoding.EncodeToString(signByte))
}

// CheckSign 检测签名
func (this *AliWebClient) CheckSign(signData, sign string) {
	signByte, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		panic(err)
	}
	s := sha1.New()
	_, err = s.Write([]byte(signData))
	if err != nil {
		panic(err)
	}
	hash := s.Sum(nil)
	err = rsa.VerifyPKCS1v15(this.PublicKey, crypto.SHA1, hash, signByte)
	if err != nil {
		panic(err)
	}
}
