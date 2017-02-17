package client

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/milkbobo/gopay/common"
	"log"
	"net/url"
	"sort"
	"strings"
)

var defaultAliAppClient *AliAppClient

type AliAppClient struct {
	PartnerID  string //合作者ID
	SellerID   string
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

func (aa *AliAppClient) Pay(charge *common.Charge) (string, error) {
	data := make(map[string]string)
	data["service"] = "mobile.securitypay.pay"
	data["partner"] = aa.PartnerID
	data["_input_charset"] = "utf-8"
	data["notify_url"] = charge.CallbackURL

	data["out_trade_no"] = charge.TradeNum
	data["subject"] = charge.Describe
	data["payment_type"] = "1"
	data["seller_id"] = aa.SellerID
	data["total_fee"] = fmt.Sprintf("%.2f", charge.MoneyFee)
	data["body"] = charge.Describe

	sign, err := aa.GenSign(data)
	if err != nil {
		return "", err
	}
	data["sign"] = sign
	data["sign_type"] = "RSA"

	var re []string
	for k, v := range data {
		re = append(re, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(re, "&"), nil
}

// GenSign 产生签名
func (aa *AliAppClient) GenSign(m map[string]string) (string, error) {
	delete(m, "sign_type")
	delete(m, "sign")
	var data []string
	for k, v := range m {
		if v == "" {
			continue
		}
		data = append(data, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(data)
	signData := strings.Join(data, "&")
	s := sha1.New()
	_, err := s.Write([]byte(signData))
	if err != nil {
		log.Println(err)
	}
	hashByte := s.Sum(nil)
	signByte, err := aa.PrivateKey.Sign(rand.Reader, hashByte, crypto.SHA1)
	if err != nil {
		return "", err
	}
	return url.QueryEscape(base64.StdEncoding.EncodeToString(signByte)), nil
}

// CheckSign 检测签名
func (aa *AliAppClient) CheckSign(data string, sign string) error {
	signByte, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	s := sha1.New()
	_, err = s.Write([]byte(data))
	if err != nil {
		return err
	}
	hash := s.Sum(nil)
	return rsa.VerifyPKCS1v15(aa.PublicKey, crypto.SHA1, hash[:], signByte)
}
