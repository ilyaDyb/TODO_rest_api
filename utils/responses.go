package utils


// MessageResponse represents a generic message response
type MessageResponse struct {
    Message string `json:"message"`
}

// ErrorResponse represents a generic error response
type ErrorResponse struct {
    Error string `json:"error"`
}