package response

type CommonResponse struct {
	IsSuccess bool        `json:"is_success"`
	Code      string      `json:"code,omitempty"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

func NewSuccessResponse(data interface{}) CommonResponse {
	return CommonResponse{
		IsSuccess: true,
		Data:      data,
	}
}

func NewErrorResponse(code string, message string) CommonResponse {
	return CommonResponse{
		IsSuccess: false,
		Code:      code,
		Message:   message,
	}
}
