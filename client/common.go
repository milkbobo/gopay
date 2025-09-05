package client

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/milkbobo/gopay/common"
	"github.com/milkbobo/gopay/util"
	"sort"
	"strings"

	"github.com/shopspring/decimal"
)

// 微信企业付款到零钱
func WachatCompanyChange(mchAppid, mchid, key string, conn *HTTPSClient, charge *common.Charge) (map[string]string, error) {
	var m = make(map[string]string)
	m["mch_appid"] = mchAppid
	m["mchid"] = mchid
	m["nonce_str"] = util.RandomStr()
	m["partner_trade_no"] = charge.TradeNum
	m["openid"] = charge.OpenID
	m["amount"] = WechatMoneyFeeToString(charge.MoneyFee)
	m["spbill_create_ip"] = util.LocalIP()
	m["desc"] = TruncatedText(charge.Describe, 32)

	// 是否验证用户名称
	if charge.CheckName {
		m["check_name"] = "FORCE_CHECK"
		m["re_user_name"] = charge.ReUserName
	} else {
		m["check_name"] = "NO_CHECK"
	}

	sign, err := WechatGenSign(key, m)
	if err != nil {
		return map[string]string{}, err
	}
	m["sign"] = sign

	// 转出xml结构
	result, err := PostWechat("https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers", m, conn)
	if err != nil {
		return map[string]string{}, err
	}

	return struct2Map(result)
}

// 微信关闭订单
func WachatCloseOrder(appid, mchid, key string, outTradeNo string) (common.WeChatQueryResult, error) {
	var m = make(map[string]string)
	m["appid"] = appid
	m["mch_id"] = mchid
	m["nonce_str"] = util.RandomStr()
	m["out_trade_no"] = outTradeNo
	m["sign_type"] = "MD5"

	sign, err := WechatGenSign(key, m)
	if err != nil {
		return common.WeChatQueryResult{}, err
	}
	m["sign"] = sign

	// 转出xml结构
	result, err := PostWechat("https://api.mch.weixin.qq.com/pay/closeorder", m, nil)
	if err != nil {
		return common.WeChatQueryResult{}, err
	}

	return result, err
}

// 微信订单查询
func WachatQueryOrder(appID, mchID, key, tradeNum string) (common.WeChatQueryResult, error) {
	var m = make(map[string]string)
	m["appid"] = appID
	m["mch_id"] = mchID
	m["out_trade_no"] = tradeNum
	m["nonce_str"] = util.RandomStr()

	sign, err := WechatGenSign(key, m)
	if err != nil {
		return common.WeChatQueryResult{}, err
	}

	m["sign"] = sign

	return PostWechat("https://api.mch.weixin.qq.com/pay/orderquery", m, nil)
}

func WechatGenSign(key string, m map[string]string) (string, error) {
	var signData []string
	for k, v := range m {
		if v != "" && k != "sign" && k != "key" {
			signData = append(signData, fmt.Sprintf("%s=%s", k, v))
		}
	}

	sort.Strings(signData)
	signStr := strings.Join(signData, "&")
	signStr = signStr + "&key=" + key

	c := md5.New()
	_, err := c.Write([]byte(signStr))
	if err != nil {
		return "", errors.New("WechatGenSign md5.Write: " + err.Error())
	}
	signByte := c.Sum(nil)
	if err != nil {
		return "", errors.New("WechatGenSign md5.Sum: " + err.Error())
	}
	return strings.ToUpper(fmt.Sprintf("%x", signByte)), nil
}

func TruncatedText(data string, length int) string {
	data = FilterTheSpecialSymbol(data)
	if len([]rune(data)) > length {
		return string([]rune(data)[:length-1])
	}
	return data
}

// 过滤特殊符号
func FilterTheSpecialSymbol(data string) string {
	// 定义转换规则
	specialSymbol := func(r rune) rune {
		if r == '`' || r == '[' || r == '~' || r == '!' || r == '@' || r == '#' || r == '$' ||
			r == '^' || r == '&' || r == '*' || r == '~' || r == '(' || r == ')' || r == '=' ||
			r == '~' || r == '|' || r == '{' || r == '}' || r == '~' || r == ':' || r == ';' ||
			r == '\'' || r == ',' || r == '\\' || r == '[' || r == ']' || r == '.' || r == '<' ||
			r == '>' || r == '/' || r == '?' || r == '~' || r == '！' || r == '@' || r == '#' ||
			r == '￥' || r == '…' || r == '&' || r == '*' || r == '（' || r == '）' || r == '—' ||
			r == '|' || r == '{' || r == '}' || r == '【' || r == '】' || r == '‘' || r == '；' ||
			r == '：' || r == '”' || r == '“' || r == '\'' || r == '"' || r == '。' || r == '，' ||
			r == '、' || r == '？' || r == '%' || r == '+' || r == '_' || r == ']' || r == '"' || r == '&' {
			return ' '
		}
		return r
	}
	data = strings.Map(specialSymbol, data)
	return strings.Replace(data, "\n", " ", -1)
}

// 对微信下订单或者查订单
func PostWechat(url string, data map[string]string, h *HTTPSClient) (common.WeChatQueryResult, error) {
	var xmlRe common.WeChatQueryResult
	buf := bytes.NewBufferString("")

	for k, v := range data {
		buf.WriteString(fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, v, k))
	}
	xmlStr := fmt.Sprintf("<xml>%s</xml>", buf.String())

	hc := new(HTTPSClient)
	if h != nil {
		hc = h
	} else {
		hc = HTTPSC
	}

	re, err := hc.PostData(url, "text/xml:charset=UTF-8", xmlStr)
	if err != nil {
		return xmlRe, errors.New("HTTPSC.PostData: " + err.Error())
	}

	err = xml.Unmarshal(re, &xmlRe)
	if err != nil {
		return xmlRe, errors.New("xml.Unmarshal: " + err.Error())
	}

	if xmlRe.ReturnCode != "SUCCESS" {
		// 通信失败
		return xmlRe, errors.New("xmlRe.ReturnMsg: " + xmlRe.ReturnMsg)
	}

	if xmlRe.ResultCode != "SUCCESS" {
		// 业务结果失败
		return xmlRe, errors.New("xmlRe.ErrCodeDes: " + xmlRe.ErrCodeDes)
	}
	return xmlRe, nil
}

// 对支付宝者查订单
func GetAlipay(url string) (common.AliWebQueryResult, error) {
	var xmlRe common.AliWebQueryResult

	re, err := HTTPSC.GetData(url)
	if err != nil {
		return xmlRe, errors.New("HTTPSC.PostData: " + err.Error())
	}
	err = xml.Unmarshal(re, &xmlRe)
	if err != nil {
		return xmlRe, errors.New("xml.Unmarshal: " + err.Error())
	}
	return xmlRe, nil
}

// 对支付宝者查订单
func GetAlipayApp(urls string) (common.AliWebAppQueryResult, error) {
	var aliPay common.AliWebAppQueryResult

	re, err := HTTPSC.GetData(urls)
	if err != nil {
		return aliPay, errors.New("HTTPSC.PostData: " + err.Error())
	}

	err = json.Unmarshal(re, &aliPay)
	if err != nil {
		panic(fmt.Sprintf("re is %v, err is %v", re, err))
	}

	return aliPay, nil
}

// ToURL
func ToURL(payUrl string, m map[string]string) string {
	var buf []string
	for k, v := range m {
		buf = append(buf, fmt.Sprintf("%s=%s", k, v))
	}
	return fmt.Sprintf("%s?%s", payUrl, strings.Join(buf, "&"))
}

// 微信金额浮点转字符串
func WechatMoneyFeeToString(moneyFee decimal.Decimal) string {
	aDecimal := moneyFee
	bDecimal := decimal.NewFromFloat(100)
	return aDecimal.Mul(bDecimal).Truncate(0).String()
}

// 支付宝金额转字符串
func AliyunMoneyFeeToString(moneyFee decimal.Decimal) string {
	return moneyFee.Truncate(2).String()
}

func struct2Map(obj interface{}) (map[string]string, error) {

	j2 := make(map[string]string)

	j1, err := json.Marshal(obj)
	if err != nil {
		return j2, err
	}

	err2 := json.Unmarshal(j1, &j2)
	return j2, err2
}
