package common

// WeChatResult 微信支付返回
type WeChatReResult struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`

	AppID      string `xml:"appid"`
	MchID      string `xml:"mch_id"`
	DeviceInfo string `xml:"device_info"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`

	TradeType string `xml:"trade_type"`
	PrepayID  string `xml:"prepay_id"`
	CodeURL   string `xml:"code_url"`
}

// WechatBaseResult 基本信息
type WechatBaseResult struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

// WechatReturnData 返回通用数据
type WechatReturnData struct {
	AppID      string `xml:"appid,emitempty"`
	MchID      string `xml:"mch_id,emitempty"`
	DeviceInfo string `xml:"device_info,emitempty"`
	NonceStr   string `xml:"nonce_str,emitempty"`
	Sign       string `xml:"sign,emitempty"`
	ResultCode string `xml:"result_code,emitempty"`
	ErrCode    string `xml:"err_code,emitempty"`
	ErrCodeDes string `xml:"err_code_des,emitempty"`
}

// WechatResultData 结果通用数据
type WechatResultData struct {
	OpenID        string `xml:"openid,emitempty"`
	IsSubscribe   string `xml:"is_subscribe,emitempty"`
	TradeType     string `xml:"trade_type,emitempty"`
	BankType      string `xml:"bank_type,emitempty"`
	FeeType       string `xml:"fee_type,emitempty"`
	TotalFee      int64  `xml:"total_fee,emitempty"`
	CashFeeType   string `xml:"cash_fee_type,emitempty"`
	CashFee       int64  `xml:"cash_fee,emitempty"`
	TransactionID string `xml:"transaction_id,emitempty"`
	OutTradeNO    string `xml:"out_trade_no,emitempty"`
	Attach        string `xml:"attach,emitempty"`
	TimeEnd       string `xml:"time_end,emitempty"`
}

type WeChatPayResult struct {
	WechatBaseResult
	WechatReturnData
	WechatResultData
}

type WeChatQueryResult struct {
	WechatBaseResult
	WechatReturnData
	WechatResultData
	TradeState     string `xml:"trade_state"`
	TradeStateDesc string `xml:"trade_state_desc"`
}
