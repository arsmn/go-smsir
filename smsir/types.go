package smsir

import (
	"fmt"
	"net/http"
)

type (
	apiResponse interface {
		isSuccessful() bool
		message() string
		messageDetail() string
	}

	BaseResponse struct {
		IsSuccessful  bool   `json:"IsSuccessful,omitempty"`
		Message       string `json:"Message,omitempty"`
		MessageDetail string `json:"MessageDetail,omitempty"`
	}

	ErrorResponse struct {
		resp          *http.Response
		Message       string
		MessageDetail string
	}
)

func (r *BaseResponse) isSuccessful() bool {
	return r.IsSuccessful
}

func (r *BaseResponse) message() string {
	return r.Message
}

func (r *BaseResponse) messageDetail() string {
	return r.MessageDetail
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %+v",
		r.resp.Request.Method, r.resp.Request.URL,
		r.resp.StatusCode, r.Message, r.MessageDetail)
}
