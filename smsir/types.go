package smsir

import (
	"fmt"
	"net/http"
)

type (
	APIResponse interface {
		IsSuccessful() bool
		Message() string
		MessageDetail() string
	}

	BaseResponse struct {
		IsResSuccessful  bool   `json:"IsSuccessful,omitempty"`
		ResMessage       string `json:"Message,omitempty"`
		ResMessageDetail string `json:"MessageDetail,omitempty"`
	}

	ErrorResponse struct {
		*http.Response
		Message       string
		MessageDetail string
	}
)

func (r *BaseResponse) IsSuccessful() bool {
	return r.IsResSuccessful
}

func (r *BaseResponse) Message() string {
	return r.ResMessage
}

func (r *BaseResponse) MessageDetail() string {
	return r.ResMessageDetail
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %+v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message, r.MessageDetail)
}
