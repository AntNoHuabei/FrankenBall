package request

type SessionListRequest struct {
	Size int `json:"size"`
	Page int `json:"page"`
}

type SessionDeleteRequest struct {
	Id string `json:"id"`
}
type SessionMessagesRequest struct {
	Session string `json:"session"`
}
