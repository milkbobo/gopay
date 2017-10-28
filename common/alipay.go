package common

// AliWebPayResult 支付宝支付结果回调
type AliWebPayResult struct {
	NotifyTime       string `json:"notify_time"`
	NotifyType       string `json:"notify_type"`
	NotifyID         string `json:"notify_id"`
	AppId            string `json:"app_id"`
	Charset          string `json:"charset"`
	Version          string `json:"version"`
	SignType         string `json:"sign_type"`
	Sign             string `json:"sign"`
	TradeNum         string `json:"trade_no"`
	OutTradeNum      string `json:"out_trade_no"`
	OutBizNo         string `json:"out_biz_no"`
	BuyerID          string `json:"buyer_id"`
	BuyerPayAmount   string `json:"buyer_pay_amount"`
	RefundFee        string `json:"refund_fee"`
	Subject          string `json:"subject"`
	PayMentType      string `json:"payment_type"`
	GmtPayMent       string `json:"gmt_payment"`
	GmtClose         string `json:"gmt_close"`
	SellerEmail      string `json:"seller_email"`
	BuyerEmail       string `json:"buyer_email"`
	SellerID         string `json:"seller_id"`
	Price            string `json:"price"`
	TotalFee         string `json:"total_fee"`
	Quantity         string `json:"quantity"`
	Body             string `json:"body"`
	Discount         string `json:"discount"`
	IsTotalFeeAdjust string `json:"is_total_fee_adjust"`
	UseCoupon        string `json:"use_coupon"`
	RefundStatus     string `json:"refund_status"`
	GmtRefund        string `json:"gmt_refund"`
	AliQueryResult
}

type FundBill struct {
	FundChannel string `json:"fundChannel"`
	Amount      string `json:"amount"`
}

type AliQueryResult struct {
	TradeNo             string     `json:"trade_no"`
	OutTradeNo          string     `json:"out_trade_no"`
	OpenID              string     `json:"open_id"`
	BuyerLogonID        string     `json:"buyer_logon_id"`
	TradeStatus         string     `json:"trade_status"`
	TotalAmount         string     `json:"total_amount"`
	ReceiptAmount       string     `json:"receipt_amount"`
	BuyerPayAmount      string     `json:"BuyerPayAmount"`
	PointAmount         string     `json:"point_amount"`
	InvoiceAmount       string     `json:"invoice_amount"`
	SendPayDate         string     `json:"send_pay_date"`
	AlipayStoreID       string     `json:"alipay_store_id"`
	StoreID             string     `json:"store_id"`
	TerminalID          string     `json:"terminal_id"`
	FundBillListStr     string     `json:"fund_bill_list"`
	FundBillList        []FundBill `json:"-"`
	StoreName           string     `json:"store_name"`
	BuyerUserID         string     `json:"buyer_user_id"`
	DiscountGoodsDetail string     `json:"discount_goods_detail"`
	IndustrySepcDetail  string     `json:"industry_sepc_detail"`
	PassbackParams      string     `json:"passback_params"`
}

type AliWebQueryResult struct {
	IsSuccess string `xml:"is_success"`
	ErrorMsg  string `xml:"error"`
	SignType  string `xml:"sign_type"`
	Sign      string `xml:"sign"`
	Response  struct {
		Trade struct {
			BuyerEmail          string `xml:"buyer_email"`
			BuyerId             string `xml:"buyer_id"`
			SellerID            string `xml:"seller_id"`
			TradeStatus         string `xml:"trade_status"`
			IsTotalFeeAdjust    string `xml:"is_total_fee_adjust"`
			OutTradeNum         string `xml:"out_trade_no"`
			Subject             string `xml:"subject"`
			FlagTradeLocked     string `xml:"flag_trade_locked"`
			Body                string `xml:"body"`
			GmtCreate           string `xml:"gmt_create"`
			GmtPayment          string `xml:"gmt_payment"`
			GmtLastModifiedTime string `xml:"gmt_last_modified_time"`
			SellerEmail         string `xml:"seller_email"`
			TotalFee            string `xml:"total_fee"`
			TradeNum            string `xml:"trade_no"`
		} `xml:"trade"`
	} `xml:"response"`
}

type AliWebAppQueryResult struct {
	AlipayTradeQueryResponse struct {
		Code                string `json:"code"`
		Msg                 string `json:"msg"`
		SubCode             string `json:"sub_code"`
		SubMsg              string `json:"sub_msg"`
		TradeNo             string `json:"trade_no"`
		OutTradeNo          string `json:"out_trade_no"`
		OpenId              string `json:"open_id"`
		BuyerLogonId        string `json:"buyer_logon_id"`
		TradeStatus         string `json:"trade_status"`
		TotalAmount         string `json:"total_amount"`
		ReceiptAmount       string `json:"receipt_amount"`
		BuyerPayAmount      string `json:"buyer_pay_amount"`
		PointAmount         string `json:"point_amount"`
		InvoiceAmount       string `json:"invoice_amount"`
		SendPayDate         string `json:"send_pay_date"`
		AlipayStoreId       string `json:"alipay_store_id"`
		StoreId             string `json:"store_id"`
		TerminalId          string `json:"terminal_id"`
		StoreName           string `json:"store_name"`
		BuyerUserId         string `json:"buyer_user_id"`
		DiscountGoodsDetail string `json:"discount_goods_detail"`
		IndustrySepcDetail  string `json:"industry_sepc_detail"`
	} `json:"alipay_trade_query_response"`
	Sign string `json:"sign"`
}
