package smsir

import (
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultUserAgent     = "go-smsir-client"
	defaultAuthHeader    = "x-sms-ir-secure-token"
	defaultTokenLifeTime = 30 * time.Minute
)

var (
	defaultBaseURL = &url.URL{
		Scheme: "https",
		Host:   "RestfulSms.com/api",
		Path:   "/",
	}
	defaultClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			MaxConnsPerHost:       100,
			MaxIdleConnsPerHost:   100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
)
