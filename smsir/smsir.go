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
	"time"

	"github.com/google/go-querystring/query"
)

type Client struct {
	client *http.Client

	baseURL   *url.URL
	userAgent string

	apiKey, secretKey string

	tokenSource   TokenSource
	tokenLifeTime time.Duration

	common service

	UserInfo     *UserInfoSerive
	SendReceive  *SendReceiveService
	Verification *VerificationService
}

type service struct {
	client *Client
}

func NewClient() *Client {
	c := &Client{
		baseURL:       defaultBaseURL,
		client:        defaultClient,
		userAgent:     defaultUserAgent,
		tokenLifeTime: defaultTokenLifeTime,
	}

	c.common.client = c

	c.UserInfo = (*UserInfoSerive)(&c.common)
	c.SendReceive = (*SendReceiveService)(&c.common)
	c.Verification = (*VerificationService)(&c.common)

	return c
}

func (c *Client) WithHttpClient(hc *http.Client) *Client {
	c.client = hc
	return c
}

func (c *Client) WithBaseURL(u *url.URL) *Client {
	c.baseURL = u
	return c
}

func (c *Client) WithUserAgent(ua string) *Client {
	c.userAgent = ua
	return c
}

func (c *Client) WithTokenSource(ts TokenSource) *Client {
	c.tokenSource = ts
	return c
}

func (c *Client) WithTokenLifeTime(lf time.Duration) *Client {
	c.tokenLifeTime = lf
	return c
}

func (c *Client) WithAuthentication(apiKey, secretKey string) *Client {
	c.apiKey = apiKey
	c.secretKey = secretKey
	c.tokenSource = ReuseTokenSource(nil, c)
	return c
}

func (c *Client) NewRequest(method, url string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.baseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.baseURL)
	}

	u, err := c.baseURL.Parse(url)
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

	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
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

func (c *Client) Do(ctx context.Context, req *http.Request, v apiResponse, auth ...bool) error {
	if ctx == nil {
		return errors.New("context must be non-nil")
	}

	shouldAuth := true
	if len(auth) > 0 {
		shouldAuth = auth[0]
	}

	if shouldAuth {
		if c.tokenSource == nil {
			return errors.New("request requires authentication and token source is not registered")
		}

		t, err := c.tokenSource.Token()
		if err != nil {
			return err
		}

		req.Header.Set(defaultAuthHeader, t.AccessToken)
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

	if !v.isSuccessful() {
		return &ErrorResponse{
			resp:          resp,
			Message:       v.message(),
			MessageDetail: v.messageDetail(),
		}
	}

	return nil
}

func (c *Client) Token() (*Token, error) {
	res, err := c.UserInfo.GetToken(context.Background(), &GetTokenRequest{
		APIKey:    c.apiKey,
		SecretKey: c.secretKey,
	})
	if err != nil {
		return nil, err
	}
	return &Token{
		AccessToken: res.Token,
		Expiry:      time.Now().Add(c.tokenLifeTime),
	}, nil
}

func addOptions(s string, opts interface{}) (string, error) {
	if opts == nil {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
