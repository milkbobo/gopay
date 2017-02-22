package gopay

import (
	"fmt"
	"github.com/milkbobo/gopay/client"
	"github.com/milkbobo/gopay/common"
	"github.com/milkbobo/gopay/constant"
	"net/http"
	"testing"
	"crypto/x509"
	"encoding/pem"
	"crypto/rsa"
)

func TestPay(t *testing.T) {
	initClient()
	initHandle()
	charge := new(common.Charge)
	charge.PayMethod = constant.ALI_WEB
	charge.MoneyFee = 1
	charge.Describe = "test pay"
	charge.TradeNum = "11111111122"
	charge.CallbackURL = "http://127.0.0.1/callback/aliappcallback"
	// charge.PayURL = "https://openapi.alipaydev.com/gateway.do"

	fdata, err := Pay(charge)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v", fdata)
}

func initClient() {

	block, _ := pem.Decode([]byte(`-----BEGIN PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAMsKJnlSm3h5xt4E
j71kGlgXueO4d9kc8JDU1og+omMlGwIfxgHruBP5QVJpvS7PUyLb7chwWWo+51qp
E5lQF497tEk2M30wcB9u4luN7/RZklc0iYqqVPbiYdpMh4yAVf5/Gkw8Ycuxsl6Z
vpLvp1MW7oWuV0gxxkDfuqO49PMXAgMBAAECgYEAqPnjlyjGtvcyKGfHcKk0u4fT
bs+A/rH1C7P2byEhaD3jQltLISIZ6pWQZZQWnDRzThmWxS+rWp7LUEpSQ0/Cqm2H
OeiDVe3Zt5Q9LxMhesqkjZRchlIQRn+Rb24wNxx3rIxi6rcVjtf7oFZ6PklSkn8U
IJkHyl5oVJk8hptT/8ECQQDpLjoW+uIwTNpKJ6+9snKB+lrE3Tc+R0ob1UzO/ild
EsDAnt2OKiF/6rzmG1+6XF7ghabaC80FGZww2JQbPTPHAkEA3ujQ0ebb76CS7TDI
7mYgUTyloL/PhNemCUiL89/IEP3KGiRSoaMhajHudfqvi2ODCpZEzDOCR5zPyDbI
If2mMQJAaxP2Sv00hzeTekAVPMhIOxXLPuHS739vMa7Wkas3NW1aJFoPpawFLCeQ
TR6+6+ZlDzdwsmp+4FutVOTvxj5pmwJAIPY2YsOLhDyvXUmYfMA3SSv5pfKXIiKt
V7QVleNidzjAGOuEGIjB2S03ANUn/imh5//efn+jZSmIBCgtofEbEQJBAL15Fi26
JJMXqKNBsVr/5dA1Kx1AjLEp6mQqfxVTdSH+AEjPYuTZFGQ9csirdqQSnWa6SrxC
dBeREcj+Yoly5h8=
-----END PRIVATE KEY-----`))

	if block == nil {
		panic("Sign private key decode error")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	block, _ = pem.Decode([]byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDLCiZ5Upt4ecbeBI+9ZBpYF7nj
uHfZHPCQ1NaIPqJjJRsCH8YB67gT+UFSab0uz1Mi2+3IcFlqPudaqROZUBePe7RJ
NjN9MHAfbuJbje/0WZJXNImKqlT24mHaTIeMgFX+fxpMPGHLsbJemb6S76dTFu6F
rldIMcZA37qjuPTzFwIDAQAB
-----END PUBLIC KEY-----`))

	if block == nil {
		panic("Sign public key decode error")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	client.InitAliWebClient(&client.AliWebClient{
		PartnerID:  "2088421851185803",
		SellerID:   "2088421851185803",
		AppID:      "2016073100134934",
		PrivateKey: privateKey.(*rsa.PrivateKey),
		PublicKey:  publicKey.(*rsa.PublicKey),
		PayURL:     "https://mapi.alipay.com/gateway.do",
	})
}

//func TestPay(t *testing.T) {
//	initClient()
//	initHandle()
//	charge := new(common.Charge)
//	charge.PayMethod = constant.WECHAT_WEB
//	charge.MoneyFee = 1
//	charge.Describe = "test pay"
//	charge.TradeNum = "11111111122"
//	charge.CallbackURL = "http://127.0.0.1/callback/aliappcallback"
//	charge.OpenID = "123"
//
//	fdata, err := Pay(charge)
//	if err != nil {
//		t.Error(err)
//	}
//	fmt.Printf("%+v", fdata)
//}
//
//func initClient() {
//	client.InitWxWebClient(&client.WechatWebClient{
//		AppID:  "wxa48fb327e06e00e5",
//		MchID:   "1437976302",
//		Key:      "hongbeibangtobenumberone20170101",
//		CallbackURL: "127.0.0.1/weixin/paycallback",
//		PayURL: "https://api.mch.weixin.qq.com/pay/unifiedorder",
//	})
//}

func initHandle() {
	http.HandleFunc("callback/aliappcallback", func(w http.ResponseWriter, r *http.Request) {
		aliResult, err := AliAppCallback(w, r)
		if err != nil {
			fmt.Println(err)
			//log.xxx
			return
		}
		selfHandler(aliResult)
	})
}

func selfHandler(i interface{}) {
}
