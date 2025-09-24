package service

import (
	"fmt"

	"new-wasa/service/models"
)

// ChatService handles business logic for chat operations
type ChatService struct {
	store ChatStore
}

// ChatStore interface defines the data access layer
type ChatStore interface {
	GetUserByToken(token string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	CreateUser(username string) (*models.User, error)
	CreateToken(userID string) (string, error)
	GetConversations(userID string) ([]models.Conversation, error)
	GetConversation(id, userID string) (*models.Conversation, error)
	CreateConversation(name, convType, userID string) (*models.Conversation, error)
	SendMessage(conversationID, userID, body string) (*models.Message, error)
	GetMessages(conversationID string) ([]models.Message, error)
	ForwardMessage(messageID, conversationID, userID string) (*models.Message, error)
	CommentMessage(messageID, userID, comment string) (*models.Comment, error)
	UncommentMessage(messageID, userID string) error
	DeleteMessage(messageID string) error
	AddToGroup(groupID, userID string) error
	LeaveGroup(groupID, userID string) error
	SetGroupName(groupID, name, userID string) error
	SetMyPhoto(userID string, photoData []byte) error
	SetGroupPhoto(groupID, userID string, photoData []byte) error
	SetMyUserName(userID, username string) error
}

// NewChatService creates a new chat service
func NewChatService(store ChatStore) *ChatService {
	return &ChatService{
		store: store,
	}
}

// Login handles user authentication
func (s *ChatService) Login(name string) (*models.LoginResponse, error) {
	// Check if user exists
	user, err := s.store.GetUserByUsername(name)
	if err != nil {
		// User doesn't exist, create new one
		user, err = s.store.CreateUser(name)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Create token for user
	token, err := s.store.CreateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	return &models.LoginResponse{
		Identifier: token,
	}, nil
}

// SetMyUserName updates user's username
func (s *ChatService) SetMyUserName(userID, username string) error {
	return s.store.SetMyUserName(userID, username)
}

// GetMyConversations retrieves user's conversations
func (s *ChatService) GetMyConversations(userID string) ([]models.Conversation, error) {
	return s.store.GetConversations(userID)
}

// GetConversation retrieves a specific conversation
func (s *ChatService) GetConversation(id, userID string) (*models.Conversation, error) {
	return s.store.GetConversation(id, userID)
}

// CreateConversation creates a new conversation
func (s *ChatService) CreateConversation(name, convType, userID string) (*models.Conversation, error) {
	return s.store.CreateConversation(name, convType, userID)
}

// SendMessage sends a message to a conversation
func (s *ChatService) SendMessage(conversationID, userID, body string) (*models.Message, error) {
	return s.store.SendMessage(conversationID, userID, body)
}

// ForwardMessage forwards a message to another conversation
func (s *ChatService) ForwardMessage(messageID, conversationID, userID string) (*models.Message, error) {
	return s.store.ForwardMessage(messageID, conversationID, userID)
}

// CommentMessage adds a comment to a message
func (s *ChatService) CommentMessage(messageID, userID, comment string) (*models.Comment, error) {
	return s.store.CommentMessage(messageID, userID, comment)
}

// UncommentMessage removes a comment from a message
func (s *ChatService) UncommentMessage(messageID, userID string) error {
	return s.store.UncommentMessage(messageID, userID)
}

// DeleteMessage deletes a message
func (s *ChatService) DeleteMessage(messageID string) error {
	return s.store.DeleteMessage(messageID)
}

// AddToGroup adds a user to a group
func (s *ChatService) AddToGroup(groupID, userID string) error {
	return s.store.AddToGroup(groupID, userID)
}

// LeaveGroup removes a user from a group
func (s *ChatService) LeaveGroup(groupID, userID string) error {
	return s.store.LeaveGroup(groupID, userID)
}

// SetGroupName updates a group's name
func (s *ChatService) SetGroupName(groupID, name, userID string) error {
	return s.store.SetGroupName(groupID, name, userID)
}

// SetMyPhoto updates user's profile photo
func (s *ChatService) SetMyPhoto(userID string, photoData []byte) error {
	return s.store.SetMyPhoto(userID, photoData)
}

// SetGroupPhoto updates a group's photo
func (s *ChatService) SetGroupPhoto(groupID, userID string, photoData []byte) error {
	return s.store.SetGroupPhoto(groupID, userID, photoData)
}
