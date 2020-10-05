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
