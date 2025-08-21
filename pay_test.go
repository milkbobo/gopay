package gopay

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"testing"

	"gopay/client"
	"gopay/common"
	"gopay/constant"
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

	// 私钥
	block, _ := pem.Decode([]byte(`-----BEGIN PRIVATE KEY-----
xxxxxxx
-----END PRIVATE KEY-----`))

	if block == nil {
		panic("Sign private key decode error")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	// 公钥
	block, _ = pem.Decode([]byte(`-----BEGIN PUBLIC KEY-----
xxxxxxxx
-----END PUBLIC KEY-----`))

	if block == nil {
		panic("Sign public key decode error")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	client.InitAliWebClient(&client.AliWebClient{
		PartnerID:  "xxxxxxxxxxxx",
		SellerID:   "xxxxxxxxxxxx",
		AppID:      "xxxxxxxxxxxx",
		PrivateKey: privateKey.(*rsa.PrivateKey),
		PublicKey:  publicKey.(*rsa.PublicKey),
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
//		AppID:       "xxxxxxxxxxxx",
//		MchID:       "xxxxxxxxxxxx",
//		Key:         "xxxxxxxxxxxx",
//		CallbackURL: "127.0.0.1/weixin/paycallback",
//		PayURL:      "https://api.mch.weixin.qq.com/pay/unifiedorder",
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
