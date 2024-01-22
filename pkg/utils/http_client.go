package utils

import (
	"Open_IM/pkg/common/config"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/guonaihong/gout"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

func HttpGetWithHeader(url string, mapvalue map[string]string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	for key, value := range mapvalue {
		req.Header.Add(key, value)
	}
	if err != nil {
		return nil, fmt.Errorf("httpGet req get error: %s", err.Error())
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("httpGet client do error: %s", err.Error())
	}
	defer resp.Body.Close()
	var body []byte
	var bodyReader io.ReadCloser
	if resp.Header.Get("Content-Encoding") == "gzip" {
		bodyReader, err = gzip.NewReader(resp.Body)
		if err != nil {
			// fallback to raw data
			bodyReader = resp.Body
		}
	} else if resp.Header.Get("Content-Encoding") == "br" {
		bodyReader = io.NopCloser(brotli.NewReader(resp.Body))
		if err != nil {
			// fallback to raw data
			bodyReader = resp.Body
		}
	} else if resp.Header.Get("Content-Encoding") == "deflate" {
		bodyReader = flate.NewReader(resp.Body)
		if err != nil {
			// fallback to raw data
			bodyReader = resp.Body
		}
	} else {
		bodyReader = resp.Body
	}
	body, err = ioutil.ReadAll(bodyReader)
	if err != nil {
		return nil, fmt.Errorf("httpGet io read error: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		//sdklog.Error("httpGet HttpReqError, op: get", "url", url, "statusCode", resp.StatusCode)
		return nil, fmt.Errorf("httpGet status code not 200, url(%s) status code %d", url, resp.StatusCode)
	}
	if body == nil || len(body) == 0 {
		return nil, fmt.Errorf("httpGet resp body nil")
	}
	return body, nil
}

func HttpGetWithHeaderWithEncode(url string, mapvalue map[string]string) ([]byte, string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	for key, value := range mapvalue {
		req.Header.Add(key, value)
	}
	if err != nil {
		return nil, "", fmt.Errorf("httpGet req get error: %s", err.Error())
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "", fmt.Errorf("httpGet client do error: %s", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, resp.Header.Get("Content-Encoding"), nil

}
func HttpGetWithHeaderWithGout(url string, mapvalue map[string]string, ptrfuncproxy string) (resultString string, code int, err error) {
	dataForma := gout.GET(url).SetHeader(mapvalue)
	if ptrfuncproxy != "" {
		dataForma = dataForma.SetProxy(ptrfuncproxy)
	}
	dataForma = dataForma.Debug(!config.Config.IsPublicEnv)
	err = dataForma.BindBody(&resultString).Code(&code).Do()
	return
}
func HttpGetWithHeaderWithProxy(url string, mapvalue map[string]string, ptrfuncproxy func(*http.Request) (*url.URL, error)) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if ptrfuncproxy != nil {
		tr.Proxy = ptrfuncproxy
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	for key, value := range mapvalue {
		req.Header.Add(key, value)
	}
	if err != nil {
		return nil, fmt.Errorf("httpGet req get error: %s", err.Error())
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("httpGet client do error: %s", err.Error())
	}
	defer resp.Body.Close()
	var body []byte
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			// 处理 gzip 解压错误
			fmt.Println("错误的解压方式")
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("httpGet io read error: %s", err.Error())
			}
		} else {
			defer reader.Close()

			uncompressedBody, err := ioutil.ReadAll(reader)
			if err != nil {
				// 处理读取解压数据错误
			}
			body = uncompressedBody
		}

	} else {
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("httpGet io read error: %s", err.Error())
		}
	}
	if resp.StatusCode != 200 {
		//sdklog.Error("httpGet HttpReqError, op: get", "url", url, "statusCode", resp.StatusCode)
		return nil, fmt.Errorf("httpGet status code not 200, url(%s) status code %d", url, resp.StatusCode)
	}
	if body == nil || len(body) == 0 {
		return nil, fmt.Errorf("httpGet resp body nil")
	}
	//sdklog.Info("httpGet success", "url", url, "resp", string(body))
	return body, nil
}
func HttpGetWithProxy(url string, ptrfuncproxy func(*http.Request) (*url.URL, error)) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if ptrfuncproxy != nil {
		tr.Proxy = ptrfuncproxy
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("httpGet req get error: %s", err.Error())
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("httpGet client do error: %s", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("httpGet io read error: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		//sdklog.Error("httpGet HttpReqError, op: get", "url", url, "statusCode", resp.StatusCode)
		return nil, fmt.Errorf("httpGet status code not 200, url(%s) status code %d", url, resp.StatusCode)
	}
	if body == nil || len(body) == 0 {
		return nil, fmt.Errorf("httpGet resp body nil")
	}
	//sdklog.Info("httpGet success", "url", url, "resp", string(body))
	return body, nil
}
func HttpGet(url string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("httpGet req get error: %s", err.Error())
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("httpGet client do error: %s", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("httpGet io read error: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("httpGet status code not 200, url(%s) status code %d", url, resp.StatusCode)
	}
	if body == nil || len(body) == 0 {
		return nil, fmt.Errorf("httpGet resp body nil")
	}
	//sdklog.Info("httpGet success", "url", url, "resp", string(body))
	return body, nil
}

func HttpPost(url string, host string, mapValueHeader map[string]string, data []byte) ([]byte, error) {
	tdetail := net.Dialer{
		KeepAlive: 10 * time.Minute,
	}
	tp := &http.Transport{
		DialContext:           tdetail.DialContext,
		ResponseHeaderTimeout: 60 * time.Second,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
	}
	client := &http.Client{Timeout: 60 * time.Second, Transport: tp}
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("httpPost req post error: %s", err.Error())
	}
	if host != "" {
		req.Host = host
	}
	if len(mapValueHeader) == 0 {
		req.Header.Add("Content-Type", "application/json")
	} else {
		for key, value := range mapValueHeader {
			req.Header.Add(key, value)
		}
	}

	//req.Header.Add("User-Agent", "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11")
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://www.apifox.cn)")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("httpPost client do error: %s", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("httpPost io read error: %s, body[%d]:%s", err.Error(), len(body), string(body))
	}
	if resp.StatusCode != 200 {
		//sdklog.Error("httpPost HttpReqError, op: post", "url", url, "statusCode", resp.StatusCode)
		return nil, fmt.Errorf("httpPost status code not 200, url(%s)  status code %d", url, resp.StatusCode)
	}
	if body == nil || len(body) == 0 {
		return nil, fmt.Errorf("httpPost resp body nil")
	}
	//sdklog.Info("httpPost success", "url", url, "resp", string(body))
	return body, nil
}

// http post参数请求
func HttpPostForm(url string, requestBody url.Values) ([]byte, error) {
	var responseByte []byte
	resp, err := http.PostForm(url, requestBody)
	if err != nil {
		return responseByte, err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return responseBody, err
	}
	return responseBody, err
}
