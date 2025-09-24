package api

import "new-wasa/service"

// NewHandler returns the HTTP handler set for the API, aligned with the template naming
func NewHandler(s *service.ChatService) *ChatHandler {
	return NewChatHandler(s)
}
