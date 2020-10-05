package smsir

import (
	"context"
	"time"
)

type SendReceiveService service

type MessageSendRequest struct {
	Messages                 []string   `json:"Messages,omitempty"`
	MobileNumbers            []string   `json:"MobileNumbers,omitempty"`
	LineNumber               *string    `json:"LineNumber,omitempty"`
	SendDateTime             *time.Time `json:"SendDateTime,omitempty"`
	CanContinueInCaseOfError *bool      `json:"CanContinueInCaseOfError,omitempty"`
}

type SentSMSLog struct {
	ID           int64  `json:"ID,omitempty"`
	MobileNumber string `json:"MobileNo,omitempty"`
}

type MessageSendResponse struct {
	BaseResponse
	Logs     []*SentSMSLog `json:"Ids,omitempty"`
	BatchKey *string       `json:"BatchKey,omitempty"`
}

func (s *SendReceiveService) SendMessage(ctx context.Context, send *MessageSendRequest) (*MessageSendResponse, error) {
	u := "MessageSend"
	req, err := s.client.Post(u, send)
	if err != nil {
		return nil, err
	}

	m := new(MessageSendResponse)
	if err := s.client.Do(ctx, req, m); err != nil {
		return nil, err
	}
	return m, nil
}

type ReportByIDOptions struct {
	ID int `url:"id,omitempty"`
}

type ReportByDateOptions struct {
	RequestedPageNumber int    `url:"RequestedPageNumber,omitempty"`
	RowsPerPage         int    `url:"RowsPerPage,omitempty"`
	ShamsiFromDate      string `url:"Shamsi_FromDate,omitempty"`
	ShamsiToDate        string `url:"Shamsi_ToDate,omitempty"`
}

type ReportByBatchKeyOptions struct {
	PageID   int `url:"pageId,omitempty"`
	BatchKey int `url:"batchKey,omitempty"`
}

type SentMessage struct {
	ID                   int64  `json:"ID,omitempty"`
	LineNumber           string `json:"LineNumber,omitempty"`
	Body                 string `json:"SMSMessageBody,omitempty"`
	MobileNumber         string `json:"MobileNo,omitempty"`
	SendDateTime         string `json:"SendDateTime,omitempty"`
	ToBeSentAt           string `json:"ToBeSentAt,omitempty"`
	NativeDeliveryStatus string `json:"NativeDeliveryStatus,omitempty"`
	Type                 string `json:"TypeOfMessage,omitempty"`
	ShamsiSendDateTime   string `json:"PersianSendDateTime,omitempty"`
	LatinReceiveDateTime string `json:"LatinReceiveDateTime,omitempty"`
}

type ReportMessagesResponse struct {
	BaseResponse
	CountOfAll int            `json:"CountOfAll,omitempty"`
	Messages   []*SentMessage `json:"Messages,omitempty"`
}

func (s *SendReceiveService) GetSentReportByDate(ctx context.Context, opts *ReportByDateOptions) (*ReportMessagesResponse, error) {
	return s.getReportByDate(ctx, "MessageSend", opts)
}

func (s *SendReceiveService) GetSentReportByBatchKey(ctx context.Context, opts *ReportByBatchKeyOptions) (*ReportMessagesResponse, error) {
	return s.getReportByDate(ctx, "MessageSend", opts)
}

func (s *SendReceiveService) GetReceiveReportByDate(ctx context.Context, opts *ReportByDateOptions) (*ReportMessagesResponse, error) {
	return s.getReportByDate(ctx, "ReceiveMessage", opts)
}

func (s *SendReceiveService) GetReceiveReportByID(ctx context.Context, opts *ReportByIDOptions) (*ReportMessagesResponse, error) {
	return s.getReportByDate(ctx, "ReceiveMessage", opts)
}

func (s *SendReceiveService) getReportByDate(ctx context.Context, u string, opts interface{}) (*ReportMessagesResponse, error) {
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.Get(u, nil)
	if err != nil {
		return nil, err
	}

	r := new(ReportMessagesResponse)
	if err := s.client.Do(ctx, req, r); err != nil {
		return nil, err
	}
	return r, nil
}

type SentMessageDetail struct {
	ID                       int64   `json:"ID,omitempty"`
	MobileNumber             *string `json:"MobileNo,omitempty"`
	SendDateTime             *string `json:"SendDateTime,omitempty"`
	DeliveryStatus           *string `json:"DeliveryStatus,omitempty"`
	Body                     *string `json:"SMSMessageBody,omitempty"`
	SendIsErronous           bool    `json:"SendIsErronous,omitempty"`
	DeliveryStatusFetchError *string `json:"DeliveryStatusFetchError,omitempty"`
	NeedsReCheck             bool    `json:"NeedsReCheck,omitempty"`
	DeliveryStateID          *int    `json:"DeliveryStateID,omitempty"`
}

type ReportByIDResponse struct {
	BaseResponse
	Message SentMessageDetail `json:"Messages,omitempty"`
}

func (s *SendReceiveService) GetSentReportByID(ctx context.Context, opts *ReportByIDOptions) (*ReportByIDResponse, error) {
	u := "MessageSend"
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.Get(u, nil)
	if err != nil {
		return nil, err
	}

	r := new(ReportByIDResponse)
	if err := s.client.Do(ctx, req, r); err != nil {
		return nil, err
	}
	return r, nil
}
