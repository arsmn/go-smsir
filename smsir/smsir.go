package smsir

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultBaseURL   = "http://RestfulSms.com/api"
	defaultUserAgent = "go-smsir"
)

type Client struct {
	client *http.Client

	BaseURL   *url.URL
	UserAgent string

	common service

	UserInfo *UserInfoSerive
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

	c.UserInfo = (*UserInfoSerive)(&c.common)

	return c
}

func (c *Client) NewRequest(method, url string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}

	u, err := c.BaseURL.Parse(url)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

func (c *Client) Get(url string, body interface{}) (*http.Request, error) {
	return c.NewRequest(http.MethodGet, url, body)
}

func (c *Client) Post(url string, body interface{}) (*http.Request, error) {
	return c.NewRequest(http.MethodPost, url, body)
}

func (c *Client) Put(url string, body interface{}) (*http.Request, error) {
	return c.NewRequest(http.MethodPut, url, body)
}

func (c *Client) Delete(url string, body interface{}) (*http.Request, error) {
	return c.NewRequest(http.MethodDelete, url, body)
}

func (c *Client) Do(ctx context.Context, req *http.Request, v APIResponse) error {
	if ctx == nil {
		return errors.New("context must be non-nil")
	}
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		return err
	}
	defer resp.Body.Close()

	if v == nil {
		v = &BaseResponse{}
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		if err != io.EOF {
			return err
		}
	}

	if !v.IsSuccessful() {
		return &ErrorResponse{
			Response:      resp,
			Message:       v.Message(),
			MessageDetail: v.MessageDetail(),
		}
	}

	return nil
}
