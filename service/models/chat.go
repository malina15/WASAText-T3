package models

import "time"

// User represents a user in the chat system
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Photo    string `json:"photo,omitempty"`
}

// Conversation represents a chat conversation
type Conversation struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // "direct" or "group"
	CreatedAt time.Time `json:"createdAt,omitempty"`
	Messages  []Message `json:"messages,omitempty"`
}

// Message represents a chat message
type Message struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversationId"`
	AuthorID       string    `json:"author"`
	Body           string    `json:"body"`
	Timestamp      time.Time `json:"timestamp"`
	Comments       []Comment `json:"comments,omitempty"`
}

// Comment represents a comment on a message
type Comment struct {
	ID        string    `json:"id"`
	MessageID string    `json:"messageId"`
	AuthorID  string    `json:"author"`
	Body      string    `json:"body"`
	Timestamp time.Time `json:"timestamp"`
}

// Group represents a group conversation
type Group struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Photo string `json:"photo,omitempty"`
}

// GroupMember represents a member of a group
type GroupMember struct {
	GroupID string `json:"groupId"`
	UserID  string `json:"userId"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Name string `json:"name"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Identifier string `json:"identifier"`
}

// SetUsernameRequest represents a set username request
type SetUsernameRequest struct {
	Username string `json:"username"`
}

// SendMessageRequest represents a send message request
type SendMessageRequest struct {
	ConversationID string `json:"conversationId"`
	Body           string `json:"body"`
}

// ForwardMessageRequest represents a forward message request
type ForwardMessageRequest struct {
	ConversationID string `json:"conversationId"`
}

// CommentMessageRequest represents a comment message request
type CommentMessageRequest struct {
	Comment string `json:"comment"`
}

// AddToGroupRequest represents an add to group request
type AddToGroupRequest struct {
	UserID string `json:"userId"`
}

// SetGroupNameRequest represents a set group name request
type SetGroupNameRequest struct {
	Name string `json:"name"`
}
