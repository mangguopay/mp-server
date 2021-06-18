package common

import (
	"a.a/cu/util"
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"
)

func PostSend(notifyUrl, body string) ([]byte, error) {
	if notifyUrl == "" {
		return []byte("缺少回调地址"), nil
	}
	client := &http.Client{
		Timeout: 1 * time.Minute,
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}
	client.Transport = tr
	httpResp, err := client.Post(notifyUrl, util.CONTENT_TYPE_XFORM, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	respText, _ := ioutil.ReadAll(httpResp.Body)
	return respText, nil
}
