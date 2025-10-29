package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"new-wasa/service/models"
)

type ChatStore interface {
	// User operations
	CreateUser(username string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	SetMyUserName(userID, username string) error
	SetMyPhoto(userID string, photoData []byte) error

	// Token operations
	CreateToken(userID string) (string, error)
	GetUserByToken(token string) (*models.User, error)
	DeleteToken(token string) error

	// Conversation operations
	CreateConversation(name, convType, userID string) (*models.Conversation, error)
	GetConversation(id, userID string) (*models.Conversation, error)
	GetConversations(userID string) ([]models.Conversation, error)

	// Message operations
	SendMessage(conversationID, userID, body string) (*models.Message, error)
	GetMessages(conversationID string) ([]models.Message, error)
	ForwardMessage(messageID, conversationID, userID string) (*models.Message, error)
	DeleteMessage(messageID string) error

	// Comment operations
	CommentMessage(messageID, userID, comment string) (*models.Comment, error)
	UncommentMessage(messageID, userID string) error

	// Group operations
	AddToGroup(groupID, userID string) error
	LeaveGroup(groupID, userID string) error
	SetGroupName(groupID, name, userID string) error
	SetGroupPhoto(groupID, userID string, photoData []byte) error
}

type SQLiteChatStore struct {
	db *sql.DB
	mu sync.RWMutex
}

func NewSQLiteChatStore(dbPath string) (*SQLiteChatStore, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=1")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	store := &SQLiteChatStore{db: db}
	if err := store.initTables(); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	return store, nil
}

func (s *SQLiteChatStore) initTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			photo TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS tokens (
			token TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS conversations (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			type TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS messages (
			id TEXT PRIMARY KEY,
			conversation_id TEXT NOT NULL,
			author_id TEXT NOT NULL,
			body TEXT NOT NULL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
			FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS comments (
			id TEXT PRIMARY KEY,
			message_id TEXT NOT NULL,
			author_id TEXT NOT NULL,
			body TEXT NOT NULL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
			FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS groups (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			photo TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS group_members (
			group_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			PRIMARY KEY (group_id, user_id),
			FOREIGN KEY (group_id) REFERENCES conversations(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
	}

	for _, query := range queries {
		if _, err := s.db.Exec(query); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}

func (s *SQLiteChatStore) CreateUser(username string) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := generateID()
	user := &models.User{
		ID:       id,
		Username: username,
	}

	_, err := s.db.Exec("INSERT INTO users (id, username) VALUES (?, ?)", user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *SQLiteChatStore) GetUserByID(id string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user := &models.User{}
	var photo sql.NullString
	err := s.db.QueryRow("SELECT id, username, photo FROM users WHERE id = ?", id).Scan(&user.ID, &user.Username, &photo)
	if err != nil {
		return nil, err
	}
	if photo.Valid {
		user.Photo = photo.String
	}

	return user, nil
}

func (s *SQLiteChatStore) GetUserByUsername(username string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user := &models.User{}
	var photo sql.NullString
	err := s.db.QueryRow("SELECT id, username, photo FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &photo)
	if err != nil {
		return nil, err
	}
	if photo.Valid {
		user.Photo = photo.String
	}

	return user, nil
}

func (s *SQLiteChatStore) UpdateUserUsername(id, username string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("UPDATE users SET username = ? WHERE id = ?", username, id)
	return err
}

func (s *SQLiteChatStore) UpdateUserPhoto(id, photo string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("UPDATE users SET photo = ? WHERE id = ?", photo, id)
	return err
}

func (s *SQLiteChatStore) CreateToken(userID string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	token := generateID()
	_, err := s.db.Exec("INSERT INTO tokens (token, user_id) VALUES (?, ?)", token, userID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *SQLiteChatStore) GetUserByToken(token string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var userID string
	err := s.db.QueryRow("SELECT user_id FROM tokens WHERE token = ?", token).Scan(&userID)
	if err != nil {
		return nil, err
	}

	// Get user directly without additional locking
	user := &models.User{}
	var photo sql.NullString
	err = s.db.QueryRow("SELECT id, username, photo FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Username, &photo)
	if err != nil {
		return nil, err
	}
	if photo.Valid {
		user.Photo = photo.String
	}

	return user, nil
}

func (s *SQLiteChatStore) DeleteToken(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("DELETE FROM tokens WHERE token = ?", token)
	return err
}

func (s *SQLiteChatStore) CreateConversation(name, convType, userID string) (*models.Conversation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := generateID()
	conv := &models.Conversation{
		ID:   id,
		Name: name,
		Type: convType,
	}

	// Create conversation
	_, err := s.db.Exec("INSERT INTO conversations (id, name, type) VALUES (?, ?, ?)", conv.ID, conv.Name, conv.Type)
	if err != nil {
		return nil, err
	}

	// Add user to conversation
	_, err = s.db.Exec("INSERT INTO group_members (group_id, user_id) VALUES (?, ?)", conv.ID, userID)
	if err != nil {
		return nil, err
	}

	return conv, nil
}

func (s *SQLiteChatStore) GetConversationByID(id string) (*models.Conversation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conv := &models.Conversation{}
	err := s.db.QueryRow("SELECT id, name, type FROM conversations WHERE id = ?", id).Scan(&conv.ID, &conv.Name, &conv.Type)
	if err != nil {
		return nil, err
	}

	return conv, nil
}

func (s *SQLiteChatStore) GetConversations(userID string) ([]models.Conversation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query(`
		SELECT DISTINCT c.id, c.name, c.type 
		FROM conversations c
		LEFT JOIN messages m ON c.id = m.conversation_id
		LEFT JOIN group_members gm ON c.id = gm.group_id
		WHERE m.author_id = ? OR gm.user_id = ?
		ORDER BY c.name
	`, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []models.Conversation
	for rows.Next() {
		var conv models.Conversation
		if err := rows.Scan(&conv.ID, &conv.Name, &conv.Type); err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}

func (s *SQLiteChatStore) GetConversation(id, userID string) (*models.Conversation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var conv models.Conversation
	err := s.db.QueryRow(`
		SELECT c.id, c.name, c.type
		FROM conversations c
		WHERE c.id = ?
	`, id).Scan(&conv.ID, &conv.Name, &conv.Type)
	if err != nil {
		return nil, err
	}

	// Get messages for this conversation
	messages, err := s.GetMessages(id)
	if err != nil {
		return nil, err
	}
	conv.Messages = messages

	return &conv, nil
}

func (s *SQLiteChatStore) CreateMessage(conversationID, authorID, body string) (*models.Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := generateID()
	message := &models.Message{
		ID:             id,
		ConversationID: conversationID,
		AuthorID:       authorID,
		Body:           body,
		Timestamp:      time.Now(),
	}

	_, err := s.db.Exec("INSERT INTO messages (id, conversation_id, author_id, body, timestamp) VALUES (?, ?, ?, ?, ?)",
		message.ID, message.ConversationID, message.AuthorID, message.Body, message.Timestamp)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (s *SQLiteChatStore) GetMessageByID(id string) (*models.Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	message := &models.Message{}
	err := s.db.QueryRow("SELECT id, conversation_id, author_id, body, timestamp FROM messages WHERE id = ?", id).
		Scan(&message.ID, &message.ConversationID, &message.AuthorID, &message.Body, &message.Timestamp)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (s *SQLiteChatStore) GetConversationMessages(conversationID string) ([]models.Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query("SELECT id, conversation_id, author_id, body, timestamp FROM messages WHERE conversation_id = ? ORDER BY timestamp", conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.ID, &message.ConversationID, &message.AuthorID, &message.Body, &message.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *SQLiteChatStore) DeleteMessage(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("DELETE FROM messages WHERE id = ?", id)
	return err
}

func (s *SQLiteChatStore) AddComment(messageID, authorID, body string) (*models.Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := generateID()
	comment := &models.Comment{
		ID:        id,
		MessageID: messageID,
		AuthorID:  authorID,
		Body:      body,
		Timestamp: time.Now(),
	}

	_, err := s.db.Exec("INSERT INTO comments (id, message_id, author_id, body, timestamp) VALUES (?, ?, ?, ?, ?)",
		comment.ID, comment.MessageID, comment.AuthorID, comment.Body, comment.Timestamp)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *SQLiteChatStore) RemoveComment(messageID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("DELETE FROM comments WHERE message_id = ?", messageID)
	return err
}

func (s *SQLiteChatStore) CreateGroup(name string) (*models.Group, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := generateID()
	group := &models.Group{
		ID:   id,
		Name: name,
	}

	_, err := s.db.Exec("INSERT INTO groups (id, name) VALUES (?, ?)", group.ID, group.Name)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (s *SQLiteChatStore) GetGroupByID(id string) (*models.Group, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	group := &models.Group{}
	err := s.db.QueryRow("SELECT id, name, photo FROM groups WHERE id = ?", id).Scan(&group.ID, &group.Name, &group.Photo)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (s *SQLiteChatStore) UpdateGroupName(id, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("UPDATE groups SET name = ? WHERE id = ?", name, id)
	return err
}

func (s *SQLiteChatStore) UpdateGroupPhoto(id, photo string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("UPDATE groups SET photo = ? WHERE id = ?", photo, id)
	return err
}

func (s *SQLiteChatStore) AddUserToGroup(groupID, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("INSERT INTO group_members (group_id, user_id) VALUES (?, ?)", groupID, userID)
	return err
}

func (s *SQLiteChatStore) RemoveUserFromGroup(groupID, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("DELETE FROM group_members WHERE group_id = ? AND user_id = ?", groupID, userID)
	return err
}

func (s *SQLiteChatStore) GetGroupMembers(groupID string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query("SELECT user_id FROM group_members WHERE group_id = ?", groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userIDs, nil
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Additional methods to match the service interface

func (s *SQLiteChatStore) SendMessage(conversationID, userID, body string) (*models.Message, error) {
	return s.CreateMessage(conversationID, userID, body)
}

func (s *SQLiteChatStore) GetMessages(conversationID string) ([]models.Message, error) {
	return s.GetConversationMessages(conversationID)
}

func (s *SQLiteChatStore) ForwardMessage(messageID, conversationID, userID string) (*models.Message, error) {
	// Get original message
	originalMsg, err := s.GetMessageByID(messageID)
	if err != nil {
		return nil, err
	}

	// Create forwarded message
	forwardedBody := "Forwarded: " + originalMsg.Body
	return s.CreateMessage(conversationID, userID, forwardedBody)
}

func (s *SQLiteChatStore) CommentMessage(messageID, userID, comment string) (*models.Comment, error) {
	return s.AddComment(messageID, userID, comment)
}

func (s *SQLiteChatStore) UncommentMessage(messageID, userID string) error {
	return s.RemoveComment(messageID)
}

func (s *SQLiteChatStore) AddToGroup(groupID, userID string) error {
	return s.AddUserToGroup(groupID, userID)
}

func (s *SQLiteChatStore) LeaveGroup(groupID, userID string) error {
	return s.RemoveUserFromGroup(groupID, userID)
}

func (s *SQLiteChatStore) SetGroupName(groupID, name, userID string) error {
	return s.UpdateGroupName(groupID, name)
}

func (s *SQLiteChatStore) SetGroupPhoto(groupID, userID string, photoData []byte) error {
	// For now, just store as base64 string
	photoStr := fmt.Sprintf("data:image/jpeg;base64,%s", string(photoData))
	return s.UpdateGroupPhoto(groupID, photoStr)
}

func (s *SQLiteChatStore) SetMyUserName(userID, username string) error {
	return s.UpdateUserUsername(userID, username)
}

func (s *SQLiteChatStore) SetMyPhoto(userID string, photoData []byte) error {
	// For now, just store as base64 string
	photoStr := fmt.Sprintf("data:image/jpeg;base64,%s", string(photoData))
	return s.UpdateUserPhoto(userID, photoStr)
}
