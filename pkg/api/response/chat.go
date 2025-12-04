package response

type ChatResponse struct {
	Content       string `json:"content"`
	ReasonContent string `json:"reason_content"`
	IndexOfDelta  int    `json:"index_of_delta"`
	RequestID     string `json:"request_id"`
	Error         error  `json:"error,omitempty"`
}
