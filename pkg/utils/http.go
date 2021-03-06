// Copyright (c) 2019 leosocy, leosocy@gmail.com
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
)

// HTTPRequestHeaders 代表发起HTTP请求时的请求头部分信息
// 其中`X-Forwarded-For`和`X-Real-Ip`可以计算出Client的公网IP
// `Via` 一般是在HTTP请求经由代理转发后增加的字段
type HTTPRequestHeaders struct {
	XForwardedFor string `json:"X-Forwarded-For"` // e.g. "1.2.3.4, 5.6.7.8, 1.2.3.4"
	XRealIP       string `json:"X-Real-Ip"`       // e.g. "1.2.3.4"
	Via           string `json:"Via"`             // e.g. "1.1 squid"
}

// RequestHeadersGetter 获取HTTP请求头的接口。
type RequestHeadersGetter interface {
	GetRequestHeaders() (headers HTTPRequestHeaders, err error)
	GetRequestHeadersUsingProxy(proxyURL string) (headers HTTPRequestHeaders, err error)
}

var (
	httpURLOfHTTPBin  = "http://httpbin.org/get?show_env=1"
	httpsURLOfHTTPBin = "https://httpbin.org/get?show_env=1"
)

// HTTPBinUtil get and parse the request header by requesting `http(s)://httpbin.org`
type HTTPBinUtil struct {
	Timeout time.Duration
}

// GetRequestHeaders implements RequestHeadersGetter.GetRequestHeaders
func (u HTTPBinUtil) GetRequestHeaders() (headers HTTPRequestHeaders, err error) {
	return u.GetRequestHeadersUsingProxy("")
}

// GetRequestHeadersUsingProxy implements RequestHeadersGetter.GetRequestHeaderUsingProxy
func (u HTTPBinUtil) GetRequestHeadersUsingProxy(proxyURL string) (headers HTTPRequestHeaders, err error) {
	if body, err := u.makeRequest(proxyURL, false); err == nil {
		return u.unmarshal(body)
	}
	return
}

func (u HTTPBinUtil) makeRequest(proxyURL string, https bool) (body []byte, err error) {
	var reqURL string
	if https {
		reqURL = httpsURLOfHTTPBin
	} else {
		reqURL = httpURLOfHTTPBin
	}
	resp, body, errs := gorequest.New().Proxy(proxyURL).Timeout(u.Timeout).Get(reqURL).EndBytes()
	if errs != nil || resp == nil || resp.StatusCode != http.StatusOK {
		return nil,
			fmt.Errorf("request %s failed, proxy [%s], https [%t]", reqURL, proxyURL, https)
	}
	return body, nil
}

func (u HTTPBinUtil) unmarshal(body []byte) (headers HTTPRequestHeaders, err error) {
	var bj map[string]interface{}
	if err = json.Unmarshal(body, &bj); err != nil {
		return
	}
	if headersBody, found := bj["headers"]; found {
		if headersBytes, err := json.Marshal(headersBody); err == nil {
			err = json.Unmarshal(headersBytes, &headers)
		}
		return
	}
	return headers, errors.New("`headers` not found in response body")
}

// ParsePublicIP resolves the public IP address of the Client based on Headers
// First parse the IP of the first record of the `X-Forwarded-For` field
// parse the `X-Real-Ip` field value if it does not exist
// If all parsing fails, return nil
func (h HTTPRequestHeaders) ParsePublicIP() (net.IP, error) {
	for _, ipStr := range strings.Split(h.XForwardedFor, ",") {
		if ip := net.ParseIP(strings.TrimSpace(ipStr)); ip != nil {
			return ip, nil
		}
	}
	if ip := net.ParseIP(strings.TrimSpace(h.XRealIP)); ip != nil {
		return ip, nil
	}
	return nil, errors.New("can't parse public ip")
}
