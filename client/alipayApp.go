package client

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/milkbobo/gopay/common"
	"net/url"
	"sort"
	"strings"
	"time"
)

var defaultAliAppClient *AliAppClient

type AliAppClient struct {
	SellerID   string //合作者ID
	AppID      string // 应用ID
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func InitAliAppClient(c *AliAppClient) {
	defaultAliAppClient = c
}

// DefaultAliAppClient 得到默认支付宝app客户端
func DefaultAliAppClient() *AliAppClient {
	return defaultAliAppClient
}

func (this *AliAppClient) Pay(charge *common.Charge)  (map[string]string, error) {
	var m = make(map[string]string)
	m["app_id"] = this.AppID
	m["method"] = "alipay.trade.app.pay"
	m["format"] = "JSON"
	m["charset"] = "utf-8"
	m["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	m["version"] = "1.0"
	m["notify_url"] = charge.CallbackURL
	m["subject"] = charge.Describe
	m["out_trade_no"] = charge.TradeNum
	m["product_code"] = "QUICK_MSECURITY_PAY"
	m["total_amount"] = fmt.Sprintf("%.2f", charge.MoneyFee)

	sign, err := this.GenSign(m)
	if err != nil {
		panic(err)
	}
	m["sign"] = sign
	m["sign_type"] = "RSA2"
	fmt.Println("sign:", sign)
	return m, nil
}

// GenSign 产生签名
func (this *AliAppClient) GenSign(m map[string]string) (string, error) {
	delete(m, "sign_type")
	delete(m, "sign")
	var data []string
	for k, v := range m {
		if v == "" {
			continue
		}
		data = append(data, fmt.Sprintf(`%s=%s`, k, v))
	}
	sort.Strings(data)
	signData := strings.Join(data, "&")
	fmt.Println(signData)
	s := sha256.New()
	_, err := s.Write([]byte(signData))
	if err != nil {
		panic(err)
	}
	hashByte := s.Sum(nil)
	signByte, err := this.PrivateKey.Sign(rand.Reader, hashByte, crypto.SHA256)
	if err != nil {
		panic(err)
	}
	return url.QueryEscape(base64.StdEncoding.EncodeToString(signByte)), nil
}

// CheckSign 检测签名
func (this *AliAppClient) CheckSign(signData, sign string) {
	signByte, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		panic(err)
	}
	s := sha256.New()
	_, err = s.Write([]byte(signData))
	if err != nil {
		panic(err)
	}
	hash := s.Sum(nil)
	err = rsa.VerifyPKCS1v15(this.PublicKey, crypto.SHA256, hash, signByte)
	if err != nil {
		panic(err)
	}
}
