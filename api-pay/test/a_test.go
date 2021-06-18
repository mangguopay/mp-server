package test

import (
	"a.a/cu/encrypt"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func p(url, md5Key string, m map[string]interface{}) {
	reqStrEnBefore := encrypt.Map2FormStr(m, md5Key, "&key=", encrypt.FIELD_ENCODED_NONE,
		[]string{}, "sign", false)
	// 全部小写
	sign := strings.ToLower(encrypt.DoMd5(reqStrEnBefore))
	m["sign"] = sign

	reqJson, _ := json.Marshal(m)
	fmt.Printf("req=[%v]", string(reqJson))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	req.Header.Set("Charset", "UTF-8")
	// 设置文件类型:
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func p2(url, md5Key, priKey string, m map[string]interface{}) {
	reqStrEnBefore := encrypt.Map2FormStr(m, md5Key, "&key=", encrypt.FIELD_ENCODED_NONE,
		[]string{}, "sign", false)
	// 全部小写
	sign := strings.ToUpper(encrypt.DoMd5(reqStrEnBefore))
	fmt.Printf("s=[%v]\n", reqStrEnBefore)
	sign2, _ := encrypt.DoRsa(encrypt.HANDLE_SIGN, encrypt.KEYTYPE_PKCS8, encrypt.RETFMT_BASE64, encrypt.HASHLENTYPE_SHA512,
		sign, priKey, encrypt.KEYFMT_PEM, encrypt.PRE_HANDLE_BYTE)
	m["sign"] = sign2

	reqJson, _ := json.Marshal(m)
	fmt.Printf("req=[%v]\n", string(reqJson))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	req.Header.Set("Charset", "UTF-8")
	// 设置文件类型:
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func TestA(t *testing.T) {
	m := map[string]interface{}{
		//"acc_no":       "2ebb2a67-502c-48cf-bbcd-807ab3aec892",
		"acc_no":       "6fdd7f0d-7de7-430f-ace7-25c4346c5e6d",
		"req_no":       strext.GetDailyId(),
		"amount":       "300000",
		"product_type": constants.ProductTypePoly,
		//"product_type": "poly",
		"notify_url": "http://127.0.0.1:5006/pay/fake_cb",
		"remark":     "http://127.0.0.1",
		"money_type": "khr",
		"role_type":  constants.RoleType_Merc,
	}

	p("http://127.0.0.1:5002/pay/prepay", "2bc4f25482faed7539f8dc1c5976ac69", m)
	//p2("http://127.0.0.1:5002/pay/prepay", "0c4398420b0ce7bae2144be520061b8c",
	//	`MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBAOrvCrJlUwUrn7Wt8h8WaaILVuqiBgWMCYRzJU0qGPcx2FPr5HgR+OGCJnOTcRRHthOTnXSzSEYWl4n0QXyvGRJ4Bwg5oXH1tc0EUDaFIDn+cxUVACkmAgoJkoKKPplPgfxVGPCPYxzmmE1q9BWum1QTJITyRrMYaD/eT/5OBpWnAgMBAAECfxuH57kAJrp0YkLoH3eFKNvUeFsGoa4EuxjKZZSlWkedj7xF6IApmwDSP69Ll/TIco3YvpruZ4nPG/MOrJ3v5wAJkByqYlau/VlE0lmPo6oIb1CORE2tsv8jeE3yHDq32+MURVe3NzrxhqMwAGmwVhh4dP6VYxtXLfTooSRk9DECQQDuNW2WX11vaE+6bje5z8TmrGdl4sWJB8PSiRPIVuyw3HSi8nAq7h8XmRyO9zWVAzgZ+RQBrGkiQHs9SWDh/i9FAkEA/Hr/Fn3HJaKMJkumMc6I9Mw+1NHKetJhMw4vehlhSI9MfZpeX7EvpmkqxXfAm7wVFpZoD3rMSQuhvihQazaZ+wJAZEKYnXIGIZ4F8LHzQoHwniZyXq/T9JkQEs0fRnNPrCEd9neUPn17GLOZEZN7Ofzg4A22Hf4zQfdh56m63WPfAQJBAO/eDrEYh/X6avoLDvbsyHMCBIN+WMa9TrkJthNjP2iHM910pkp1dNa9vXPjpKqQUtylwnFKcgDHuz+E33osYrcCQQDEPfM4tIFvHw4tcEfqVqWtZg+Q9c79nGUCKT4Dg9/UxTTUdNTLCs3o5bXFiI2evnLEA9oTE1WEculQuvczyfVU`, m)
}

func TestB(t *testing.T) {
	m := map[string]interface{}{
		"acc_no": "b89976dc-f039-4ad0-8dbe-b8125e8a93b7",
		"req_no": "2020032512115098649170",
	}
	//p("http://127.0.0.1:5005/pay/query", "06f6075863ed4a99c4fb006775491e02db8f08f06229cbafb8d45bd6a9d98c7c", m)
	p2("http://127.0.0.1:5005/pay/query", "0c4398420b0ce7bae2144be520061b8c",
		`MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBAOrvCrJlUwUrn7Wt8h8WaaILVuqiBgWMCYRzJU0qGPcx2FPr5HgR+OGCJnOTcRRHthOTnXSzSEYWl4n0QXyvGRJ4Bwg5oXH1tc0EUDaFIDn+cxUVACkmAgoJkoKKPplPgfxVGPCPYxzmmE1q9BWum1QTJITyRrMYaD/eT/5OBpWnAgMBAAECfxuH57kAJrp0YkLoH3eFKNvUeFsGoa4EuxjKZZSlWkedj7xF6IApmwDSP69Ll/TIco3YvpruZ4nPG/MOrJ3v5wAJkByqYlau/VlE0lmPo6oIb1CORE2tsv8jeE3yHDq32+MURVe3NzrxhqMwAGmwVhh4dP6VYxtXLfTooSRk9DECQQDuNW2WX11vaE+6bje5z8TmrGdl4sWJB8PSiRPIVuyw3HSi8nAq7h8XmRyO9zWVAzgZ+RQBrGkiQHs9SWDh/i9FAkEA/Hr/Fn3HJaKMJkumMc6I9Mw+1NHKetJhMw4vehlhSI9MfZpeX7EvpmkqxXfAm7wVFpZoD3rMSQuhvihQazaZ+wJAZEKYnXIGIZ4F8LHzQoHwniZyXq/T9JkQEs0fRnNPrCEd9neUPn17GLOZEZN7Ofzg4A22Hf4zQfdh56m63WPfAQJBAO/eDrEYh/X6avoLDvbsyHMCBIN+WMa9TrkJthNjP2iHM910pkp1dNa9vXPjpKqQUtylwnFKcgDHuz+E33osYrcCQQDEPfM4tIFvHw4tcEfqVqWtZg+Q9c79nGUCKT4Dg9/UxTTUdNTLCs3o5bXFiI2evnLEA9oTE1WEculQuvczyfVU`, m)
}

func TestC(t *testing.T) {
	m := map[string]interface{}{
		//"acc_no": "2ebb2a67-502c-48cf-bbcd-807ab3aec892",
		"acc_no": "6fdd7f0d-7de7-430f-ace7-25c4346c5e6d",
		"req_no": strext.GetDailyId(),
		"amount": "10000",
		//"transfer_type": constants.ProductType_withdraw,
		"transfer_type": constants.ProductTypeTransfer,
		"card_no":       "123",
		"card_acc":      "123",
		"product_type":  constants.ProductTypeTransfer,
		"role_type":     constants.RoleType_Merc,
		"notify_url":    "http://127.0.0.1:5006/pay/fake_cb",
	}
	p("http://127.0.0.1:5005/pay/transfer", "2bc4f25482faed7539f8dc1c5976ac69", m)
	//p2("http://127.0.0.1:5005/pay/transfer", "0c4398420b0ce7bae2144be520061b8c",
	//	`MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBAOrvCrJlUwUrn7Wt8h8WaaILVuqiBgWMCYRzJU0qGPcx2FPr5HgR+OGCJnOTcRRHthOTnXSzSEYWl4n0QXyvGRJ4Bwg5oXH1tc0EUDaFIDn+cxUVACkmAgoJkoKKPplPgfxVGPCPYxzmmE1q9BWum1QTJITyRrMYaD/eT/5OBpWnAgMBAAECfxuH57kAJrp0YkLoH3eFKNvUeFsGoa4EuxjKZZSlWkedj7xF6IApmwDSP69Ll/TIco3YvpruZ4nPG/MOrJ3v5wAJkByqYlau/VlE0lmPo6oIb1CORE2tsv8jeE3yHDq32+MURVe3NzrxhqMwAGmwVhh4dP6VYxtXLfTooSRk9DECQQDuNW2WX11vaE+6bje5z8TmrGdl4sWJB8PSiRPIVuyw3HSi8nAq7h8XmRyO9zWVAzgZ+RQBrGkiQHs9SWDh/i9FAkEA/Hr/Fn3HJaKMJkumMc6I9Mw+1NHKetJhMw4vehlhSI9MfZpeX7EvpmkqxXfAm7wVFpZoD3rMSQuhvihQazaZ+wJAZEKYnXIGIZ4F8LHzQoHwniZyXq/T9JkQEs0fRnNPrCEd9neUPn17GLOZEZN7Ofzg4A22Hf4zQfdh56m63WPfAQJBAO/eDrEYh/X6avoLDvbsyHMCBIN+WMa9TrkJthNjP2iHM910pkp1dNa9vXPjpKqQUtylwnFKcgDHuz+E33osYrcCQQDEPfM4tIFvHw4tcEfqVqWtZg+Q9c79nGUCKT4Dg9/UxTTUdNTLCs3o5bXFiI2evnLEA9oTE1WEculQuvczyfVU`, m)
}

func TestD(t *testing.T) {
	m := map[string]interface{}{
		//"acc_no": "2ebb2a67-502c-48cf-bbcd-807ab3aec892",
		"acc_no": "b89976dc-f039-4ad0-8dbe-b8125e8a93b7",
		"req_no": "2020032513533725149170",
	}
	//p("http://127.0.0.1:5005/pay/transfer_query", "06f6075863ed4a99c4fb006775491e02db8f08f06229cbafb8d45bd6a9d98c7c", m)
	p2("http://127.0.0.1:5005/pay/transfer_query", "0c4398420b0ce7bae2144be520061b8c",
		`MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBAOrvCrJlUwUrn7Wt8h8WaaILVuqiBgWMCYRzJU0qGPcx2FPr5HgR+OGCJnOTcRRHthOTnXSzSEYWl4n0QXyvGRJ4Bwg5oXH1tc0EUDaFIDn+cxUVACkmAgoJkoKKPplPgfxVGPCPYxzmmE1q9BWum1QTJITyRrMYaD/eT/5OBpWnAgMBAAECfxuH57kAJrp0YkLoH3eFKNvUeFsGoa4EuxjKZZSlWkedj7xF6IApmwDSP69Ll/TIco3YvpruZ4nPG/MOrJ3v5wAJkByqYlau/VlE0lmPo6oIb1CORE2tsv8jeE3yHDq32+MURVe3NzrxhqMwAGmwVhh4dP6VYxtXLfTooSRk9DECQQDuNW2WX11vaE+6bje5z8TmrGdl4sWJB8PSiRPIVuyw3HSi8nAq7h8XmRyO9zWVAzgZ+RQBrGkiQHs9SWDh/i9FAkEA/Hr/Fn3HJaKMJkumMc6I9Mw+1NHKetJhMw4vehlhSI9MfZpeX7EvpmkqxXfAm7wVFpZoD3rMSQuhvihQazaZ+wJAZEKYnXIGIZ4F8LHzQoHwniZyXq/T9JkQEs0fRnNPrCEd9neUPn17GLOZEZN7Ofzg4A22Hf4zQfdh56m63WPfAQJBAO/eDrEYh/X6avoLDvbsyHMCBIN+WMa9TrkJthNjP2iHM910pkp1dNa9vXPjpKqQUtylwnFKcgDHuz+E33osYrcCQQDEPfM4tIFvHw4tcEfqVqWtZg+Q9c79nGUCKT4Dg9/UxTTUdNTLCs3o5bXFiI2evnLEA9oTE1WEculQuvczyfVU`, m)
	fmt.Println("")
}
