package gopay

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/milkbobo/gopay/client"
	"github.com/milkbobo/gopay/common"
	"github.com/milkbobo/gopay/util"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
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
	if m["sign_type"] != "RSA" {
		panic("签名类型未知")
	}

	client.DefaultAliAppClient().CheckSign(signData, m["sign"])

	var aliPay common.AliWebPayResult
	err := util.MapStringToStruct(m, &aliPay)
	if err != nil {
		w.Write([]byte("error"))
		panic(err)
	}

	w.Write([]byte("success"))
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
	//body := cb.Ctx.Input.RequestBody
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//log.Error(string(body))
		returnCode = "FAIL"
		returnMsg = "Bodyerror"
		panic(err)
	}
	err = xml.Unmarshal(body, &reXML)
	if err != nil {
		//log.Error(err, string(body))
		returnMsg = "参数错误"
		returnCode = "FAIL"
		panic(err)
	}

	if reXML.ReturnCode != "SUCCESS" {
		//log.Error(reXML)
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
	WeChatPaySign := m["sign"]
	mySign, err := client.DefaultWechatWebClient().GenSign(m)
	if err != nil {
		panic(err)
	}

	if mySign != WeChatPaySign {
		panic(errors.New("签名交易错误"))
	}

	returnCode = "SUCCESS"
	return &reXML, nil
}
