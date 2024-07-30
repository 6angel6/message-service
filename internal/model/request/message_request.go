package request

// MessageRequest represents a request to create a message
// @Description Represents a request to create a message
type MessageRequest struct {
	Content string `json:"content"`
}
