package pay

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"a.a/cu/strext"
	"a.a/mp-server/common/ss_func"
	"a.a/mp-server/merchant-mock/conf"
	"a.a/mp-server/merchant-mock/myrsa"
)

type NotifyPaymentData struct {
	//Code string `json:"code"`
	//Msg  string `json:"msg"`
	NotifyType string `json:"notify_type"`

	Amount       string `json:"amount"`
	AppId        string `json:"app_id"`
	CurrencyType string `json:"currency_type"`
	OrderNo      string `json:"order_no"`
	OutOrderNo   string `json:"out_order_no"`
	OrderStatus  string `json:"order_status"`
	PayAccount   string `json:"pay_account"`
	PayTime      string `json:"pay_time"`
	Subject      string `json:"subject"`
	SignType     string `json:"sign_type"`
	Timestamp    string `json:"timestamp"`
	Sign         string `json:"sign"`
}

// 支付-异步通知验证异
func NotifyPaymentVerify(body []byte) (NotifyPaymentData, error) {
	var ret NotifyPaymentData
	reqData := strext.FormString2Map(string(body))
	//reqData := strext.Json2Map(body)
	if reqData == nil {
		return ret, errors.New("请求包体不是json格式,body:" + string(body))
	}

	reqSign, exists := reqData[SignKey]
	if !exists {
		log.Printf("返回数据中没有sign字段")
		return ret, errors.New("返回数据中没有sign字段")
	}

	// 将参数排序并拼接成字符串
	dataRespStr := ss_func.ParamsMapToString(reqData, SignKey)

	log.Printf("dataRespStr:%v", dataRespStr)

	// rsa2验签
	verifyErr := myrsa.RSA2Verify(dataRespStr, fmt.Sprintf("%s", reqSign), conf.PlatformPublicKey)
	if verifyErr != nil {
		log.Printf("verifyErr:%v", verifyErr)
		return ret, verifyErr
	}

	log.Printf("verifyOk")
	body = []byte(strext.ToJson(reqData))
	if jerr := json.Unmarshal(body, &ret); jerr != nil {
		log.Printf("接口返回数据,json解码到结构体失败, err:%v \n", jerr)
		return ret, jerr
	}

	return ret, nil
}

type NotifyRefundData struct {
	//Code string `json:"code"`
	//Msg  string `json:"msg"`

	AppId      string `json:"app_id"`
	SignType   string `json:"sign_type"`
	Timestamp  string `json:"timestamp"`
	Sign       string `json:"sign"`
	NotifyType string `json:"notify_type"`

	OrderNo      string `json:"order_no"`
	OutOrderNo   string `json:"out_order_no"`
	TransAmount  string `json:"trans_amount"`
	RefundNo     string `json:"refund_no"`
	OutRefundNo  string `json:"out_refund_no"`
	RefundAmount string `json:"refund_amount"`
	CurrencyType string `json:"currency_type"`
	RefundStatus string `json:"refund_status"`
	RefundTime   string `json:"refund_time"`
}

// 退款-异步通知验证异
func NotifyRefundVerify(body []byte) (NotifyRefundData, error) {
	var ret NotifyRefundData
	reqData := strext.FormString2Map(string(body))
	//reqData := strext.Json2Map(body)
	if reqData == nil {
		return ret, errors.New("请求包体不是json格式,body:" + string(body))
	}

	reqSign, exists := reqData[SignKey]
	if !exists {
		log.Printf("返回数据中没有sign字段")
		return ret, errors.New("返回数据中没有sign字段")
	}

	// 将参数排序并拼接成字符串
	dataRespStr := ss_func.ParamsMapToString(reqData, SignKey)

	log.Printf("dataRespStr:%v", dataRespStr)

	// rsa2验签
	verifyErr := myrsa.RSA2Verify(dataRespStr, fmt.Sprintf("%s", reqSign), conf.PlatformPublicKey)
	if verifyErr != nil {
		log.Printf("verifyErr:%v", verifyErr)
		return ret, verifyErr
	}

	log.Printf("verifyOk")
	body = []byte(strext.ToJson(reqData))
	if jerr := json.Unmarshal(body, &ret); jerr != nil {
		log.Printf("接口返回数据,json解码到结构体失败, err:%v \n", jerr)
		return ret, jerr
	}

	return ret, nil
}

type NotifyTransferData struct {
	//Code string `json:"code"`
	//Msg  string `json:"msg"`
	NotifyType string `json:"notify_type"`

	Amount       string `json:"amount"`
	AppId        string `json:"app_id"`
	CurrencyType string `json:"currency_type"`

	TransferNo     string `json:"transfer_no"`
	OutTransferNo  string `json:"out_transfer_no"`
	TransferStatus string `json:"transfer_status"`
	TransferTime   string `json:"transfer_time"`

	Subject   string `json:"subject"`
	SignType  string `json:"sign_type"`
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
}

// 企业付款-异步通知验证异
func NotifyTransferVerify(body []byte) (NotifyTransferData, error) {
	var ret NotifyTransferData
	reqData := strext.FormString2Map(string(body))
	//reqData := strext.Json2Map(body)
	if reqData == nil {
		return ret, errors.New("请求包体不是json格式,body:" + string(body))
	}

	reqSign, exists := reqData[SignKey]
	if !exists {
		log.Printf("返回数据中没有sign字段")
		return ret, errors.New("返回数据中没有sign字段")
	}

	// 将参数排序并拼接成字符串
	dataRespStr := ss_func.ParamsMapToString(reqData, SignKey)

	log.Printf("dataRespStr:%v", dataRespStr)

	// rsa2验签
	verifyErr := myrsa.RSA2Verify(dataRespStr, fmt.Sprintf("%s", reqSign), conf.PlatformPublicKey)
	if verifyErr != nil {
		log.Printf("verifyErr:%v", verifyErr)
		return ret, verifyErr
	}

	log.Printf("verifyOk")
	body = []byte(strext.ToJson(reqData))
	if jerr := json.Unmarshal(body, &ret); jerr != nil {
		log.Printf("接口返回数据,json解码到结构体失败, err:%v \n", jerr)
		return ret, jerr
	}

	return ret, nil
}
