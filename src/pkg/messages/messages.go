package messages

type MessageRequest struct {
	Content string `json:"content"`
}

type MessageResponse struct {
	Content string `json:"content"`
}