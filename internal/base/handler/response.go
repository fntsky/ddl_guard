package handler

import (
	apperrors "github.com/fntsky/ddl_guard/internal/errors"
)

type resp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

type Response = resp

func NewRespBodyData(code int, message string, data any) *resp {
	return &resp{
		Code:    string(apperrors.CodeSuccess),
		Message: message,
		Data:    data,
	}
}

func NewErrorResp(err *apperrors.AppError, requestID string) *resp {
	return &resp{
		Code:      string(err.Code),
		Message:   err.Message,
		RequestID: requestID,
	}
}
