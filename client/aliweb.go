package client

import (
	"bytes"
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

var aliWebClient *AliWebClient

// AliWebClient 支付宝网页支付
type AliWebClient struct {
	PartnerID   string          // 支付宝合作身份ID
	SellerID    string          // 卖家支付宝用户号
	AppID       string          // 支付宝分配给开发者的应用ID ps: 查询订单用
	CallbackURL string          // 回调接口
	PrivateKey  *rsa.PrivateKey // 私钥
	PublicKey   *rsa.PublicKey  // 公钥
	PayURL      string          // 支付网管地址
}

func InitAliWebClient(c *AliWebClient) {
	aliWebClient = c
}

// DefaultAliWebClient 默认支付宝网页支付客户端
func DefaultAliWebClient() *AliWebClient {
	return aliWebClient
}

// Pay 实现支付接口
func (ac *AliWebClient) Pay(charge *common.Charge) (string, error) {
	var m = make(map[string]string)
	m["service"] = "create_direct_pay_by_user"
	m["partner"] = ac.PartnerID
	m["_input_charset"] = "UTF-8"
	m["notify_url"] = ac.CallbackURL
	m["return_url"] = charge.ReturnURL
	m["out_trade_no"] = charge.TradeNum
	m["subject"] = charge.Describe
	m["total_fee"] = fmt.Sprintf("%.2f", charge.MoneyFee)
	m["seller_id"] = ac.SellerID
	m["payment_type"] = "1"
	m["show_url"] = charge.ShowURL
	sign, err := ac.GenSign(m)
	if err != nil {
		return "", err
	}
	m["sign"] = sign //
	fmt.Println("sign:", sign)
	m["sign_type"] = "RSA"

	return ac.ToHTML(m), nil
}

// GenSign 产生签名
func (ac *AliWebClient) GenSign(m map[string]string) (string, error) {
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
	s := sha1.New()
	_, err := s.Write([]byte(signData))
	if err != nil {
		log.Println(err)
	}
	hashByte := s.Sum(nil)
	signByte, err := ac.PrivateKey.Sign(rand.Reader, hashByte, crypto.SHA1)
	if err != nil {
		return "", err
	}
	return url.QueryEscape(base64.StdEncoding.EncodeToString(signByte)), nil
}

// ToHTML 转换form表单
func (ac *AliWebClient) ToHTML(m map[string]string) string {
	buf := bytes.NewBufferString("")
	for k, v := range m {
		buf.WriteString(fmt.Sprintf(`<input type="hidden" name="%s" value="%s">`, k, v))
	}
	formatStr :=
		`<html>
    <meta http-equiv=Content-Type content="text/html;charset=utf-8">
    <body>
        <form id="paysubmit" name="paysubmit" action="%s" method = "GET">
        %s
        <input type="submit" value="ok" style="display:none">
        </form>
        <script>
         (function(){
             document.forms["paysubmit"].submit();
         })();
        </script>
    </body>
</html>`
	return fmt.Sprintf(formatStr, ac.PayURL, buf.String())
}

// ToURL
func (ac *AliWebClient) ToURL(m map[string]string) string {
	var buf []string
	for k, v := range m {
		buf = append(buf, fmt.Sprintf("%s=%s", k, v))
	}
	return fmt.Sprintf("%s?%s", ac.PayURL, strings.Join(buf, "&"))
}

// CheckSign 检测签名
func (ac *AliWebClient) CheckSign(signData, sign string) error {
	signByte, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	s := sha1.New()
	_, err = s.Write([]byte(signData))
	if err != nil {
		return err
	}
	hash := s.Sum(nil)
	return rsa.VerifyPKCS1v15(ac.PublicKey, crypto.SHA1, hash, signByte)
}
