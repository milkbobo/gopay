package common

// WeChatResult 微信支付返回
type WeChatReResult struct {
	//ReturnCode string `xml:"return_code"`
	//ReturnMsg  string `xml:"return_msg"`

	//AppID      string `xml:"appid"`
	//MchID      string `xml:"mch_id"`
	//DeviceInfo string `xml:"device_info"`
	//NonceStr   string `xml:"nonce_str"`
	//Sign       string `xml:"sign"`
	//ResultCode string `xml:"result_code"`
	//ErrCode    string `xml:"err_code"`
	//ErrCodeDes string `xml:"err_code_des"`

	//TradeType string `xml:"trade_type"`
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
	AppID      string `xml:"appid,omitempty"`
	MchID      string `xml:"mch_id,omitempty"`
	DeviceInfo string `xml:"device_info,omitempty"`
	NonceStr   string `xml:"nonce_str,omitempty"`
	Sign       string `xml:"sign,omitempty"`
	ResultCode string `xml:"result_code,omitempty"`
	ErrCode    string `xml:"err_code,omitempty"`
	ErrCodeDes string `xml:"err_code_des,omitempty"`
}

// WechatResultData 结果通用数据
type WechatResultData struct {
	OpenID        string `xml:"openid,omitempty"`
	IsSubscribe   string `xml:"is_subscribe,omitempty"`
	TradeType     string `xml:"trade_type,omitempty"`
	BankType      string `xml:"bank_type,omitempty"`
	FeeType       string `xml:"fee_type,omitempty"`
	TotalFee      int64  `xml:"total_fee,omitempty"`
	CashFeeType   string `xml:"cash_fee_type,omitempty"`
	CashFee       int64  `xml:"cash_fee,omitempty"`
	TransactionID string `xml:"transaction_id,omitempty"`
	OutTradeNO    string `xml:"out_trade_no,omitempty"`
	Attach        string `xml:"attach,omitempty"`
	TimeEnd       string `xml:"time_end,omitempty"`
}

type WeChatPayResult struct {
	WechatBaseResult
	WechatReturnData
	WechatResultData
}

type WeChatQueryResult struct {
	WechatBaseResult
	WeChatReResult
	WechatReturnData
	WechatResultData
	TradeState     string `xml:"trade_state"`
	TradeStateDesc string `xml:"trade_state_desc"`
}
