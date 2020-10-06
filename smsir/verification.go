package smsir

import "context"

type VerificationService service

type VerificationCodeRequest struct {
	Code         string `json:"Code,omitempty"`
	MobileNumber string `json:"MobileNumber,omitempty"`
}

type VerificationCodeResponse struct {
	BaseResponse
	VerificationCodeID int `json:"VerificationCodeId,omitempty"`
}

func (s *VerificationService) VerificationCode(ctx context.Context, send *VerificationCodeRequest) (*VerificationCodeResponse, error) {
	u := "VerificationCode"
	req, err := s.client.Post(u, send)
	if err != nil {
		return nil, err
	}

	v := new(VerificationCodeResponse)
	if err := s.client.Do(ctx, req, v); err != nil {
		return nil, err
	}
	return v, nil
}

type UltraFastParameter struct {
	Key   string `json:"Parameter,omitempty"`
	Value string `json:"ParameterValue,omitempty"`
}

type UltraFastSendRequest struct {
	Mobile     string               `json:"Mobile,omitempty"`
	TemplateID string               `json:"TemplateId,omitempty"`
	Parameters []UltraFastParameter `json:"ParameterArray,omitempty"`
}

type ultraFastSendResponse struct {
	BaseResponse
	VerificationCodeID float64 `json:"VerificationCodeId,omitempty"`
}

type UltraFastSendResponse struct {
	BaseResponse
	VerificationCodeID int `json:"VerificationCodeId,omitempty"`
}

func (s *VerificationService) UltraFastSend(ctx context.Context, send *UltraFastSendRequest) (*UltraFastSendResponse, error) {
	u := "UltraFastSend"
	req, err := s.client.Post(u, send)
	if err != nil {
		return nil, err
	}

	f := new(ultraFastSendResponse)
	if err := s.client.Do(ctx, req, f); err != nil {
		return nil, err
	}
	return &UltraFastSendResponse{
		BaseResponse:       f.BaseResponse,
		VerificationCodeID: int(f.VerificationCodeID),
	}, nil
}
