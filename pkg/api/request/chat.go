package request

type ChatRequest struct {
	Message   string `json:"message"`
	Session   string `json:"session"`
	RequestId string `json:"request_id"`
}
