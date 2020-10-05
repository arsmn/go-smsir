package smsir

import "context"

type UserInfoSerive service

type GetTokenRequest struct {
	APIKey    string `json:"UserApiKey,omitempty"`
	SecretKey string `json:"SecretKey,omitempty"`
}

type GetTokenResponse struct {
	BaseResponse
	Token string `json:"TokenKey,omitempty"`
}

func (s *UserInfoSerive) GetToken(ctx context.Context, token *GetTokenRequest) (*GetTokenResponse, error) {
	u := "token"
	req, err := s.client.Post(u, token)
	if err != nil {
		return nil, err
	}

	t := new(GetTokenResponse)
	if err := s.client.Do(ctx, req, t); err != nil {
		return nil, err
	}
	return t, nil
}

type CreditResponse struct {
	BaseResponse
	Credit float64 `json:"Credit,omitempty"`
}

func (s *UserInfoSerive) Credit(ctx context.Context) (*CreditResponse, error) {
	u := "Credit"
	req, err := s.client.Get(u, nil)
	if err != nil {
		return nil, err
	}

	c := new(CreditResponse)
	if err := s.client.Do(ctx, req, c); err != nil {
		return nil, err
	}
	return c, nil
}

type SMSLine struct {
	ID         int `json:"ID,omitempty"`
	LineNumber int `json:"LineNumber,omitempty"`
}

type SMSLineResponse struct {
	BaseResponse
	SMSLines []*SMSLine `json:"SMSLines,omitempty"`
}

func (s *UserInfoSerive) SMSLine(ctx context.Context) (*SMSLineResponse, error) {
	u := "SMSLine"
	req, err := s.client.Get(u, nil)
	if err != nil {
		return nil, err
	}

	l := new(SMSLineResponse)
	if err := s.client.Do(ctx, req, l); err != nil {
		return nil, err
	}
	return l, nil
}
