package gopay

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/milkbobo/gopay/client"
	"github.com/milkbobo/gopay/common"
	"github.com/milkbobo/gopay/util"
	"encoding/json"
)

func AliWebCallback(w http.ResponseWriter, r *http.Request) (*common.AliWebPayResult, error) {
	var m = make(map[string]string)
	var signSlice []string
	r.ParseForm()
	for k, v := range r.Form {
		// k不会有多个值的情况
		m[k] = v[0]
		if k == "sign" || k == "sign_type" {
			continue
		}
		signSlice = append(signSlice, fmt.Sprintf("%s=%s", k, v[0]))
	}

	sort.Strings(signSlice)
	signData := strings.Join(signSlice, "&")
	if m["sign_type"] != "RSA" {
		//错误日志
		panic("签名类型未知")
	}

	client.DefaultAliWebClient().CheckSign(signData, m["sign"])

	var aliPay common.AliWebPayResult
	err := util.MapStringToStruct(m, &aliPay)
	if err != nil {
		w.Write([]byte("error"))
		panic(err)
	}

	w.Write([]byte("success"))
	return &aliPay, nil
}

// 支付宝app支付回调
func AliAppCallback(w http.ResponseWriter, r *http.Request) (*common.AliWebPayResult, error) {
	var result string
	defer func() {
		w.Write([]byte(result))
	}()

	var m = make(map[string]string)
	var signSlice []string
	r.ParseForm()
	for k, v := range r.Form {
		m[k] = v[0]
		if k == "sign" || k == "sign_type" {
			continue
		}
		signSlice = append(signSlice, fmt.Sprintf("%s=%s", k, v[0]))
	}
	sort.Strings(signSlice)
	signData := strings.Join(signSlice, "&")
	if m["sign_type"] != "RSA2" {
		result = "error"
		panic("签名类型未知")
	}

	client.DefaultAliAppClient().CheckSign(signData, m["sign"])

	mByte, err := json.Marshal(m)
	if err != nil {
		result = "error"
		panic(err)
	}

	var aliPay common.AliWebPayResult
	err = json.Unmarshal(mByte, &aliPay)
	if err != nil {
		result = "error"
		panic(fmt.Sprintf("m is %v, err is %v", m, err))
	}
	result = "success"
	return &aliPay, nil
}

// WeChatCallback 微信支付
func WeChatCallback(w http.ResponseWriter, r *http.Request) (*common.WeChatPayResult, error) {
	var returnCode = "FAIL"
	var returnMsg = ""
	defer func() {
		formatStr := `<xml><return_code><![CDATA[%s]]></return_code>
                  <return_msg>![CDATA[%s]]</return_msg></xml>`
		returnBody := fmt.Sprintf(formatStr, returnCode, returnMsg)
		w.Write([]byte(returnBody))
	}()
	var reXML common.WeChatPayResult
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		returnCode = "FAIL"
		returnMsg = "Bodyerror"
		return &reXML, errors.New(returnCode + ":" + returnMsg)
	}
	err = xml.Unmarshal(body, &reXML)
	if err != nil {
		returnCode = "FAIL"
		returnMsg = "参数错误"
		return &reXML, errors.New(returnCode + ":" + returnMsg)
	}

	if reXML.ReturnCode != "SUCCESS" {
		returnCode = "FAIL"
		return &reXML, errors.New(reXML.ReturnCode)
	}
	m := util.XmlToMap(body)

	var signData []string
	for k, v := range m {
		if k == "sign" {
			continue
		}
		signData = append(signData, fmt.Sprintf("%v=%v", k, v))
	}

	key := client.DefaultWechatAppClient().Key

	mySign, err := client.WechatGenSign(key, m)
	if err != nil {
		return &reXML, err
	}

	if mySign != m["sign"] {
		panic(errors.New("签名交易错误"))
	}

	returnCode = "SUCCESS"
	return &reXML, nil
}

func WeChatWebCallback(w http.ResponseWriter, r *http.Request) (*common.WeChatPayResult, error) {
	return WeChatCallback(w, r)
}

func WeChatAppCallback(w http.ResponseWriter, r *http.Request) (*common.WeChatPayResult, error) {
	return WeChatCallback(w, r)
}
