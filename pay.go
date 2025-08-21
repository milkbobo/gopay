package gopay

import (
	"errors"
	"gopay/client"
	"gopay/common"
	"gopay/constant"
)

// 用户下单支付接口
func Pay(charge *common.Charge) (map[string]string, error) {
	err := checkCharge(charge)
	if err != nil {
		return map[string]string{}, err
	}

	ct := getPayType(charge.PayMethod)
	re, err := ct.Pay(charge)
	return re, err
}

// 付款给用户接口
func PayToClient(charge *common.Charge) (map[string]string, error) {
	err := checkCharge(charge)
	if err != nil {
		return nil, err
	}
	ct := getPayType(charge.PayMethod)
	re, err := ct.PayToClient(charge)
	return re, err
}

// 验证支付内容
func checkCharge(charge *common.Charge) error {
	if charge.PayMethod <= 0 {
		return errors.New("PayMethod不能少于等于0")
	}
	if charge.MoneyFee <= 0 {
		return errors.New("MoneyFee不能少于等于0")
	}
	return nil
}

// getPayType 得到需要支付的类型
func getPayType(payMethod int64) common.PayClient {
	//如果使用余额支付
	switch payMethod {
	case constant.ALI_WEB:
		return client.DefaultAliWebClient()
	case constant.ALI_APP:
		return client.DefaultAliAppClient()
	case constant.WECHAT_WEB:
		return client.DefaultWechatWebClient()
	case constant.WECHAT_APP:
		return client.DefaultWechatAppClient()
	case constant.WECHAT_MINI_PROGRAM:
		return client.DefaultWechatMiniProgramClient()
	}
	return nil
}
