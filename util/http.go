package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var httpClient *http.Client

type HttpResult struct {
	Code float64
	Msg  string
	Data string
}

func init() {
	tr := &http.Transport{ //解决x509: certificate signed by unknown authority
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient = &http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}
}

// 发送GET请求
// url:请求地址
// response:请求返回的内容
func Get(url string) (response string, error error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		error = err
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()

	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			error = err
			return
		}
	}
	response = result.String()
	return
}

// 发送POST请求
// url:请求地址，data:POST请求提交的数据,contentType:请求体格式，如：application/json
// content:请求放回的内容
func Post(url string, data interface{}, contentType string) (content string, error error) {
	jsonStr, _ := json.Marshal(data)
	body := bytes.NewBuffer(jsonStr)
	req, err := http.NewRequest("POST", url, body)
	req.Header.Add("content-type", contentType)
	req.Header.Set("charset", "UTF-8")
	if err != nil {
		error = err
		return
	}
	defer func() {
		if err := req.Body.Close(); err != nil {
			return
		}
	}()

	resp, err := httpClient.Do(req)
	if err != nil {
		error = err
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()

	result, _ := ioutil.ReadAll(resp.Body)
	content = string(result)
	return
}
