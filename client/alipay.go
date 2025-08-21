package client

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gopay/common"
	"hash"
	"net/url"
	"sort"
	"strings"
	"time"
)

type AliPayClient struct {
	AppID string // 应用ID

	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey

	RSAType string // RSA or RSA2
}

func (this *AliPayClient) PayToClient(charge *common.Charge) (map[string]string, error) {
	var result = make(map[string]string)
	var m = make(map[string]string)
	var bizContent = make(map[string]string)
	m["app_id"] = this.AppID
	m["method"] = "alipay.fund.trans.toaccount.transfer"
	m["format"] = "JSON"
	m["charset"] = "utf-8"
	m["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	m["version"] = "1.0"
	m["sign_type"] = this.RSAType

	bizContent["out_biz_no"] = charge.TradeNum
	bizContent["amount"] = AliyunMoneyFeeToString(charge.MoneyFee)
	bizContent["payee_account"] = charge.AliAccount
	bizContent["payee_type"] = charge.AliAccountType

	bizContent["remark"] = TruncatedText(charge.Describe, 32)

	bizContentJson, err := json.Marshal(bizContent)
	if err != nil {
		return result, errors.New("json.Marshal: " + err.Error())
	}
	m["biz_content"] = string(bizContentJson)

	m["sign"] = this.GenSign(m)

	requestUrl := fmt.Sprintf("%s?%s", "https://openapi.alipay.com/gateway.do", this.ToURL(m))

	var resp map[string]interface{}
	bytes, err := HTTPSC.GetData(requestUrl)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return result, err
	}

	result, ok := resp["alipay_fund_trans_toaccount_transfer_response"].(map[string]string)
	if !ok {
		return result, errors.New(fmt.Sprintf("返回结果错误:%s", resp))
	}

	return result, nil
}

// GenSign 产生签名
func (this *AliPayClient) GenSign(m map[string]string) string {
	var data []string

	for k, v := range m {
		if v != "" && k != "sign" {
			data = append(data, fmt.Sprintf(`%s=%s`, k, v))
		}
	}
	sort.Strings(data)
	signData := strings.Join(data, "&")

	s := this.getHash(this.RSAType)

	_, err := s.Write([]byte(signData))
	if err != nil {
		panic(err)
	}
	hashByte := s.Sum(nil)
	signByte, err := this.PrivateKey.Sign(rand.Reader, hashByte, crypto.SHA256)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(signByte)
}

// CheckSign 检测签名
func (this *AliPayClient) CheckSign(signData, sign string) {
	signByte, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		panic(err)
	}
	s := this.getHash(this.RSAType)
	_, err = s.Write([]byte(signData))
	if err != nil {
		panic(err)
	}
	hashByte := s.Sum(nil)
	err = rsa.VerifyPKCS1v15(this.PublicKey, this.getCrypto(), hashByte, signByte)
	if err != nil {
		panic(err)
	}
}

// ToURL
func (this *AliPayClient) ToURL(m map[string]string) string {
	var buf []string
	for k, v := range m {
		buf = append(buf, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
	}
	return strings.Join(buf, "&")
}

func (this *AliPayClient) getRsa() string {
	if this.RSAType == "" {
		this.RSAType = "RSA"
	}

	return this.RSAType
}

func (this *AliPayClient) getCrypto() crypto.Hash {
	if this.RSAType == "RSA2" {
		return crypto.SHA256
	}
	return crypto.SHA1

}

func (this *AliPayClient) getHash(rasType string) hash.Hash {
	if rasType == "RSA2" {
		return sha256.New()
	}
	return sha1.New()
}
