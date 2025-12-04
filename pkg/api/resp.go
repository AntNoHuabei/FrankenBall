package api

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Success(d any) Response {
	return Response{
		Code:    200,
		Message: "success",
		Data:    d,
	}
}

func Fail(message string) Response {
	return Response{
		Code:    500,
		Message: message,
	}
}
