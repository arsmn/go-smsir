package smsir

import (
	"net/http"
	"net/url"
)

const (
	defaultBaseURL   = "http://RestfulSms.com/api"
	defaultUserAgent = "go-smsir"
)

type Client struct {
	client *http.Client

	common service

	BaseURL   *url.URL
	UserAgent string
}

type service struct {
	client *Client
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: defaultUserAgent}
	c.common.client = c

	return c
}
