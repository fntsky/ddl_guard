package handler

type resp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewRespBodyData(code int, message string, data any) *resp {
	return &resp{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
