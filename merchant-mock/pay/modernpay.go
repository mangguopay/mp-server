package pay

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_func"

	cryptorand "crypto/rand"

	"a.a/cu/encrypt"
	"a.a/mp-server/merchant-mock/conf"
	"a.a/mp-server/merchant-mock/dao"
	"a.a/mp-server/merchant-mock/myrsa"
)

const (
	CodeSuccess = "0"    // 成功
	SignKey     = "sign" // 验证字段

	SubCodeSuccess = "SUCCESS" // 成功

	DefaultLang = "zh_CN" // 默认语言(zh_CN, km_KH, en_US)
)

type PrepayResponse struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code"`
	SubMsg  string `json:"sub_msg"`

	OrderNo    string `json:"order_no"`
	OutOrderNo string `json:"out_order_no"`
	QrCode     string `json:"qr_code"`
}

// modernpay下单接口
func ModernPayPreOrder(order *dao.Order, qrcodeContent string) (PrepayResponse, error) {
	var ret PrepayResponse

	data := make(map[string]interface{})
	data["app_id"] = conf.AppId
	data["sign_type"] = "RSA2"
	data["timestamp"] = fmt.Sprintf("%v", time.Now().Unix())
	data["out_order_no"] = order.OrderSn
	data["currency_type"] = order.CurrencyType
	data["amount"] = fmt.Sprintf("%v", order.Amount)
	data["notify_url"] = conf.NotifyUrl
	data["subject"] = order.Title
	data["attach"] = "123456"
	data["time_expire"] = fmt.Sprintf("%v", time.Now().Add(time.Minute*10).Unix()) // 绝对的超时时间
	data["trade_type"] = order.TradeType
	data["lang"] = DefaultLang

	rawData, _ := json.Marshal(data)
	log.Printf("rawData:%v \n", string(rawData))

	if qrcodeContent != "" {
		// 商家扫码用户时才有
		data["payment_code"] = qrcodeContent
	}

	dataReqStr := ss_func.ParamsMapToString(data, SignKey)
	log.Printf("dataReqStr:%v \n", dataReqStr)

	reqSign, err := myrsa.RSA2Sign(dataReqStr, conf.SelfPrivateKey)
	if err != nil {
		log.Printf("RSA2Sign-err:%v \n", err)
		return ret, err
	}
	log.Printf("reqSign:%v \n", reqSign)

	data[SignKey] = reqSign

	reqData, err := json.Marshal(data)
	if err != nil {
		log.Printf("json-err:%v \n", err)
		return ret, err
	}

	log.Printf("requestData:%v \n", string(reqData))

	//return ret, errors.New("停止请求")

	respData, err := HttpPostJson(conf.PrepayUrl, reqData)
	if err != nil {
		log.Printf("HttpPostJson-err:%v \n", err)
		return ret, err
	}

	log.Printf("responseData:%v \n", string(respData))

	respMap := make(map[string]interface{})
	jerr := json.Unmarshal(respData, &respMap)
	if jerr != nil {
		log.Printf("接口返回数据,json解码失败, err:%v \n", jerr)
		return ret, jerr
	}

	respSign, exists := respMap[SignKey]
	if !exists {
		log.Printf("返回数据中没有sign字段")
		return ret, errors.New("返回数据中没有sign字段")
	}

	// 将参数排序并拼接成字符串
	dataRespStr := ss_func.ParamsMapToString(respMap, SignKey)

	log.Printf("dataRespStr:%v", dataRespStr)

	// rsa2验签
	verifyErr := myrsa.RSA2Verify(dataRespStr, fmt.Sprintf("%s", respSign), conf.PlatformPublicKey)
	if verifyErr != nil {
		log.Printf("verifyErr:%v", verifyErr)
		return ret, verifyErr
	}

	log.Printf("verifyOk")

	if jerr := json.Unmarshal(respData, &ret); jerr != nil {
		log.Printf("接口返回数据,json解码到结构体失败, err:%v \n", jerr)
		return ret, jerr
	}

	if ret.Code != CodeSuccess {
		log.Printf("接口返回code不是正确状态,code:%v, msg:%v\n", ret.Code, ret.Msg)
		return ret, errors.New(ret.Msg)
	}

	if ret.SubCode != SubCodeSuccess {
		log.Printf("接口返回sub_code不是正确状态,SubCode:%v, SubMsg:%v\n", ret.SubCode, ret.SubMsg)
		return ret, errors.New(ret.SubMsg)
	}

	return ret, nil
}

type PayResponse struct {
	Code       string `json:"code"`
	Msg        string `json:"msg"`
	OrderNo    string `json:"order_no"`
	OutOrderNo string `json:"out_order_no"`
}

// 加密支付密码
func EncryptPassword(originPwd string, nonstr string) string {
	// md5( md5(sha1(pwd).toLocaleLowerCase() + PASSWORDSIGN ).toLocaleLowerCase() + noStr).toLocaleLowerCase()
	slat := "sa5d6g728ttg$%43JASHGFUIa72"

	tmp := strings.ToLower(encrypt.DoShaXXX(originPwd, encrypt.HASHLENTYPE_SHA1))
	tmp = strings.ToLower(encrypt.DoMd5(tmp + slat))
	tmp = strings.ToLower(encrypt.DoMd5(tmp + nonstr))

	return tmp
}

// Random generate string
func GetRandomString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	cryptorand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

type QueryPayResponse struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code"`
	SubMsg  string `json:"sub_msg"`

	OutOrderNo   string `json:"out_order_no"`
	OrderNo      string `json:"order_no"`
	OrderStatus  string `json:"order_status"`
	Amount       string `json:"amount"`
	CurrencyType string `json:"currency_type"`
	CreateTime   string `json:"create_time"`
	PayTime      string `json:"pay_time"`
	Subject      string `json:"subject"`
	Remark       string `json:"remark"`
	Rate         string `json:"rate"`
	Fee          string `json:"fee"`
}

//支付结果查询接口
func QueryPay(payOrderSn, outOrderNo string) (*QueryPayResponse, error) {
	var ret *QueryPayResponse
	data := make(map[string]interface{})
	data["app_id"] = conf.AppId
	data["sign_type"] = "RSA2"
	data["timestamp"] = fmt.Sprintf("%v", time.Now().Unix())
	data["lang"] = DefaultLang

	if payOrderSn != "" {
		data["order_no"] = payOrderSn
	}

	if outOrderNo != "" {
		data["out_order_no"] = outOrderNo
	}

	//rawData, _ := json.Marshal(data)
	//log.Printf("rawData:%v \n", string(rawData))

	dataReqStr := ss_func.ParamsMapToString(data, SignKey)
	log.Printf("dataReqStr:%v \n", dataReqStr)

	reqSign, err := myrsa.RSA2Sign(dataReqStr, conf.SelfPrivateKey)
	if err != nil {
		log.Printf("RSA2Sign-err:%v \n", err)
		return nil, err
	}
	data[SignKey] = reqSign

	reqData, err := json.Marshal(data)
	if err != nil {
		log.Printf("json-err:%v \n", err)
		return nil, err
	}

	log.Printf("requestData:%v \n", string(reqData))

	// return ret, errors.New("停止请求")

	respData, err := HttpPostJson(conf.QueryPay, reqData)
	if err != nil {
		log.Printf("HttpPostJson-err:%v \n", err)
		return nil, err
	}

	log.Printf("responseData:%v \n", string(respData))

	respMap := make(map[string]interface{})
	jerr := json.Unmarshal(respData, &respMap)
	if jerr != nil {
		log.Printf("接口返回数据,json解码失败, err:%v \n", jerr)
		return nil, jerr
	}

	respSign, exists := respMap[SignKey]
	if !exists {
		log.Printf("返回数据中没有sign字段")
		return nil, errors.New("返回数据中没有sign字段")
	}

	// 将参数排序并拼接成字符串
	dataRespStr := ss_func.ParamsMapToString(respMap, SignKey)

	log.Printf("dataRespStr:%v", dataRespStr)

	// rsa2验签
	verifyErr := myrsa.RSA2Verify(dataRespStr, fmt.Sprintf("%s", respSign), conf.PlatformPublicKey)
	if verifyErr != nil {
		log.Printf("verifyErr:%v", verifyErr)
		return nil, verifyErr
	}

	log.Printf("verifyOk")

	if jerr := json.Unmarshal(respData, &ret); jerr != nil {
		log.Printf("接口返回数据,json解码到结构体失败, err:%v \n", jerr)
		return nil, jerr
	}

	if ret.Code != CodeSuccess {
		log.Printf("接口返回code不是正确状态,code:%v, msg:%v\n", ret.Code, ret.Msg)
		return nil, errors.New(ret.Msg)
	}

	if ret.SubCode != SubCodeSuccess {
		log.Printf("接口返回sub_code不是正确状态,SubCode:%v, SubMsg:%v\n", ret.SubCode, ret.SubMsg)
		return ret, errors.New(ret.SubMsg)
	}

	return ret, nil
}

func HttpPostJson(hostUrl string, jsonBytes []byte) ([]byte, error) {
	request, err := http.NewRequest("POST", hostUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		log.Printf("NewRequest-err:%v \n", err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("Client-do-err:%v \n", err)
		return nil, err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ReadAll-body-err:%v \n", err)
		return nil, err
	}

	return respBytes, nil
}

type RefundOrderResponse struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code"`
	SubMsg  string `json:"sub_msg"`

	OrderNo      string `json:"order_no"`
	OutOrderNo   string `json:"out_order_no"`
	RefundNo     string `json:"refund_no"`
	OutRefundNo  string `json:"out_refund_no"`
	RefundStatus string `json:"refund_status"`
	RefundAmount string `json:"refund_amount"`
}

// modernpay退款接口
func ModernPayRefundOrder(payOrderSn, outOrderNo, amount, outrefundNo string) (RefundOrderResponse, error) {
	var ret RefundOrderResponse

	data := make(map[string]interface{})
	data["app_id"] = conf.AppId
	data["sign_type"] = "RSA2"
	data["timestamp"] = fmt.Sprintf("%v", time.Now().Unix())
	data["refund_amount"] = amount
	data["out_refund_no"] = outrefundNo
	data["refund_reason"] = "买错了"
	data["lang"] = DefaultLang
	data["notify_url"] = conf.NotifyUrl

	if payOrderSn != "" {
		data["order_no"] = payOrderSn
	}

	if outOrderNo != "" {
		data["out_order_no"] = outOrderNo
	}

	//rawData, _ := json.Marshal(data)
	//log.Printf("rawData:%v \n", string(rawData))

	dataReqStr := ss_func.ParamsMapToString(data, SignKey)
	log.Printf("dataReqStr:%v \n", dataReqStr)

	reqSign, err := myrsa.RSA2Sign(dataReqStr, conf.SelfPrivateKey)
	if err != nil {
		log.Printf("RSA2Sign-err:%v \n", err)
		return ret, err
	}
	log.Printf("reqSign:%v \n", reqSign)

	data[SignKey] = reqSign

	reqData, err := json.Marshal(data)
	if err != nil {
		log.Printf("json-err:%v \n", err)
		return ret, err
	}

	log.Printf("requestData:%v \n", string(reqData))

	//return ret, errors.New("停止请求")

	respData, err := HttpPostJson(conf.RefundUrl, reqData)
	if err != nil {
		log.Printf("HttpPostJson-err:%v \n", err)
		return ret, err
	}

	log.Printf("responseData:%v \n", string(respData))

	respMap := make(map[string]interface{})
	jerr := json.Unmarshal(respData, &respMap)
	if jerr != nil {
		log.Printf("接口返回数据,json解码失败, err:%v \n", jerr)
		return ret, jerr
	}

	respSign, exists := respMap[SignKey]
	if !exists {
		log.Printf("返回数据中没有sign字段")
		return ret, errors.New("返回数据中没有sign字段")
	}

	// 将参数排序并拼接成字符串
	dataRespStr := ss_func.ParamsMapToString(respMap, SignKey)

	log.Printf("dataRespStr:%v", dataRespStr)

	// rsa2验签
	verifyErr := myrsa.RSA2Verify(dataRespStr, fmt.Sprintf("%s", respSign), conf.PlatformPublicKey)
	if verifyErr != nil {
		log.Printf("verifyErr:%v", verifyErr)
		return ret, verifyErr
	}

	log.Printf("verifyOk")

	if jerr := json.Unmarshal(respData, &ret); jerr != nil {
		log.Printf("接口返回数据,json解码到结构体失败, err:%v \n", jerr)
		return ret, jerr
	}

	if ret.Code != CodeSuccess {
		log.Printf("接口返回code不是正确状态,code:%v, msg:%v\n", ret.Code, ret.Msg)
		return ret, errors.New(ret.Msg)
	}

	if ret.SubCode != SubCodeSuccess {
		log.Printf("接口返回sub_code不是正确状态,SubCode:%v, SubMsg:%v\n", ret.SubCode, ret.SubMsg)
		return ret, errors.New(ret.SubMsg)
	}

	return ret, nil
}

type RefundQueryOrderResponse struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code"`
	SubMsg  string `json:"sub_msg"`

	OrderNo    string `json:"order_no"`
	OutOrderNo string `json:"out_order_no"`
	QrCode     string `json:"qr_code"`
}

// modernpay退款订单查询接口
func ModernPayRefundQueryOrder(refundNo, outRefundNo string) (RefundQueryOrderResponse, error) {
	var ret RefundQueryOrderResponse

	data := make(map[string]interface{})
	data["app_id"] = conf.AppId
	data["sign_type"] = "RSA2"
	data["timestamp"] = fmt.Sprintf("%v", time.Now().Unix())
	data["lang"] = DefaultLang

	if refundNo != "" {
		data["refund_no"] = refundNo
	}

	if outRefundNo != "" {
		data["out_refund_no"] = outRefundNo
	}

	//rawData, _ := json.Marshal(data)
	//log.Printf("rawData:%v \n", string(rawData))

	dataReqStr := ss_func.ParamsMapToString(data, SignKey)
	log.Printf("dataReqStr:%v \n", dataReqStr)

	reqSign, err := myrsa.RSA2Sign(dataReqStr, conf.SelfPrivateKey)
	if err != nil {
		log.Printf("RSA2Sign-err:%v \n", err)
		return ret, err
	}
	//log.Printf("reqSign:%v \n", reqSign)

	data[SignKey] = reqSign

	reqData, err := json.Marshal(data)
	if err != nil {
		log.Printf("json-err:%v \n", err)
		return ret, err
	}

	log.Printf("requestData:%v \n", string(reqData))

	//return ret, errors.New("停止请求")

	respData, err := HttpPostJson(conf.RefundQueryUrl, reqData)
	if err != nil {
		log.Printf("HttpPostJson-err:%v \n", err)
		return ret, err
	}

	log.Printf("responseData:%v \n", string(respData))

	respMap := make(map[string]interface{})
	jerr := json.Unmarshal(respData, &respMap)
	if jerr != nil {
		log.Printf("接口返回数据,json解码失败, err:%v \n", jerr)
		return ret, jerr
	}

	respSign, exists := respMap[SignKey]
	if !exists {
		log.Printf("返回数据中没有sign字段")
		return ret, errors.New("返回数据中没有sign字段")
	}

	// 将参数排序并拼接成字符串
	dataRespStr := ss_func.ParamsMapToString(respMap, SignKey)

	log.Printf("dataRespStr:%v", dataRespStr)

	// rsa2验签
	verifyErr := myrsa.RSA2Verify(dataRespStr, fmt.Sprintf("%s", respSign), conf.PlatformPublicKey)
	if verifyErr != nil {
		log.Printf("verifyErr:%v", verifyErr)
		return ret, verifyErr
	}

	log.Printf("verifyOk")

	if jerr := json.Unmarshal(respData, &ret); jerr != nil {
		log.Printf("接口返回数据,json解码到结构体失败, err:%v \n", jerr)
		return ret, jerr
	}

	if ret.Code != CodeSuccess {
		log.Printf("接口返回code不是正确状态,code:%v, msg:%v\n", ret.Code, ret.Msg)
		return ret, errors.New(ret.Msg)
	}

	if ret.SubCode != SubCodeSuccess {
		log.Printf("接口返回sub_code不是正确状态,SubCode:%v, SubMsg:%v\n", ret.SubCode, ret.SubMsg)
		return ret, errors.New(ret.SubMsg)
	}

	return ret, nil
}

type EnterpriseTransferResponse struct {
	Code          string `json:"code"`
	Msg           string `json:"msg"`
	SubCode       string `json:"sub_code"`
	SubMsg        string `json:"sub_msg"`
	Attach        string `json:"attach"`
	TransferNo    string `json:"transfer_no"`
	OutTransferNo string `json:"out_transfer_no"`
	Amount        string `json:"amount"`
	CurrencyType  string `json:"currency_type"`
	Status        string `json:"status"`
}

// modernpay企业付款接口
func ModernPayEnterpriseTransfer(enterpriseTransferData *dao.Transfer) (*EnterpriseTransferResponse, error) {
	var ret *EnterpriseTransferResponse
	data := make(map[string]interface{})
	data["app_id"] = conf.AppId
	data["sign_type"] = "RSA2"
	data["timestamp"] = fmt.Sprintf("%v", time.Now().Unix())
	data["lang"] = DefaultLang
	data["notify_url"] = conf.NotifyUrl
	data["out_transfer_no"] = enterpriseTransferData.OutTransferNo
	data["amount"] = fmt.Sprintf("%v", enterpriseTransferData.Amount)
	data["currency_type"] = enterpriseTransferData.CurrencyType
	data["remark"] = enterpriseTransferData.Remark

	if enterpriseTransferData.CountryCode != "" {
		data["country_code"] = enterpriseTransferData.CountryCode
	}

	if enterpriseTransferData.PayeePhone != "" {
		data["payee_phone"] = enterpriseTransferData.PayeePhone
	}

	if enterpriseTransferData.PayeeEmail != "" {
		data["payee_email"] = enterpriseTransferData.PayeeEmail
	}

	rawData, _ := json.Marshal(data)
	log.Printf("rawData:%v \n", string(rawData))

	dataReqStr := ss_func.ParamsMapToString(data, SignKey)
	log.Printf("dataReqStr:%v \n", dataReqStr)

	reqSign, err := myrsa.RSA2Sign(dataReqStr, conf.SelfPrivateKey)
	if err != nil {
		log.Printf("RSA2Sign-err:%v \n", err)
		return nil, err
	}
	data[SignKey] = reqSign

	reqData, err := json.Marshal(data)
	if err != nil {
		log.Printf("json-err:%v \n", err)
		return nil, err
	}

	log.Printf("requestData:%v \n", string(reqData))

	// return ret, errors.New("停止请求")

	respData, err := HttpPostJson(conf.EnterpriseTransferUrl, reqData)
	if err != nil {
		log.Printf("HttpPostJson-err:%v \n", err)
		return nil, err
	}

	log.Printf("responseData:%v \n", string(respData))

	respMap := make(map[string]interface{})
	jerr := json.Unmarshal(respData, &respMap)
	if jerr != nil {
		log.Printf("接口返回数据,json解码失败, err:%v \n", jerr)
		return nil, jerr
	}

	respSign, exists := respMap[SignKey]
	if !exists {
		log.Printf("返回数据中没有sign字段")
		return nil, errors.New("返回数据中没有sign字段")
	}

	// 将参数排序并拼接成字符串
	dataRespStr := ss_func.ParamsMapToString(respMap, SignKey)

	log.Printf("dataRespStr:%v", dataRespStr)

	// rsa2验签
	verifyErr := myrsa.RSA2Verify(dataRespStr, fmt.Sprintf("%s", respSign), conf.PlatformPublicKey)
	if verifyErr != nil {
		log.Printf("verifyErr:%v", verifyErr)
		return nil, verifyErr
	}

	log.Printf("verifyOk")

	if jerr := json.Unmarshal(respData, &ret); jerr != nil {
		log.Printf("接口返回数据,json解码到结构体失败, err:%v \n", jerr)
		return nil, jerr
	}

	if ret.Code != CodeSuccess {
		log.Printf("接口返回code不是正确状态,code:%v, msg:%v\n", ret.Code, ret.Msg)
		return nil, errors.New(ret.Msg)
	}

	if ret.SubCode != SubCodeSuccess {
		log.Printf("接口返回sub_code不是正确状态,SubCode:%v, SubMsg:%v\n", ret.SubCode, ret.SubMsg)
		return ret, errors.New(ret.SubMsg)
	}

	return ret, nil
}

type EnterpriseTransferQueryResponse struct {
	Code           string `json:"code"`
	Msg            string `json:"msg"`
	SubCode        string `json:"sub_code"`
	SubMsg         string `json:"sub_msg"`
	Amount         string `json:"amount"`
	CurrencyType   string `json:"currency_type"`
	TransferNo     string `json:"transfer_no"`
	OutTransferNo  string `json:"out_transfer_no"`
	TransferStatus string `json:"transfer_status"`
	TransferTime   string `json:"transfer_time"`
	WrongReason    string `json:"wrong_reason"`
}

// modernpay企业付款接口
func ModernPayEnterpriseTransferQuery(outTransferNo, transferNo string) (*EnterpriseTransferQueryResponse, error) {
	var ret *EnterpriseTransferQueryResponse
	data := make(map[string]interface{})
	data["app_id"] = conf.AppId
	data["sign_type"] = "RSA2"
	data["timestamp"] = fmt.Sprintf("%v", time.Now().Unix())
	data["lang"] = constants.LangZhCN

	if outTransferNo != "" {
		data["out_transfer_no"] = outTransferNo
	}

	if transferNo != "" {
		data["transfer_no"] = transferNo
	}
	rawData, _ := json.Marshal(data)
	log.Printf("rawData:%v \n", string(rawData))

	dataReqStr := ss_func.ParamsMapToString(data, SignKey)
	log.Printf("dataReqStr:%v \n", dataReqStr)

	reqSign, err := myrsa.RSA2Sign(dataReqStr, conf.SelfPrivateKey)
	if err != nil {
		log.Printf("RSA2Sign-err:%v \n", err)
		return nil, err
	}
	data[SignKey] = reqSign

	reqData, err := json.Marshal(data)
	if err != nil {
		log.Printf("json-err:%v \n", err)
		return nil, err
	}

	log.Printf("requestData:%v \n", string(reqData))

	// return ret, errors.New("停止请求")

	respData, err := HttpPostJson(conf.EnterpriseTransferQueryUrl, reqData)
	if err != nil {
		log.Printf("HttpPostJson-err:%v \n", err)
		return nil, err
	}

	log.Printf("responseData:%v \n", string(respData))

	respMap := make(map[string]interface{})
	jerr := json.Unmarshal(respData, &respMap)
	if jerr != nil {
		log.Printf("接口返回数据,json解码失败, err:%v \n", jerr)
		return nil, jerr
	}

	respSign, exists := respMap[SignKey]
	if !exists {
		log.Printf("返回数据中没有sign字段")
		return nil, errors.New("返回数据中没有sign字段")
	}

	// 将参数排序并拼接成字符串
	dataRespStr := ss_func.ParamsMapToString(respMap, SignKey)

	log.Printf("dataRespStr:%v", dataRespStr)

	// rsa2验签
	verifyErr := myrsa.RSA2Verify(dataRespStr, fmt.Sprintf("%s", respSign), conf.PlatformPublicKey)
	if verifyErr != nil {
		log.Printf("verifyErr:%v", verifyErr)
		return nil, verifyErr
	}

	log.Printf("verifyOk")

	if jerr := json.Unmarshal(respData, &ret); jerr != nil {
		log.Printf("接口返回数据,json解码到结构体失败, err:%v \n", jerr)
		return nil, jerr
	}

	if ret.Code != CodeSuccess {
		log.Printf("接口返回code不是正确状态,code:%v, msg:%v\n", ret.Code, ret.Msg)
		return nil, errors.New(ret.Msg)
	}

	if ret.SubCode != SubCodeSuccess {
		log.Printf("接口返回sub_code不是正确状态,SubCode:%v, SubMsg:%v\n", ret.SubCode, ret.SubMsg)
		return ret, errors.New(ret.SubMsg)
	}

	return ret, nil
}
