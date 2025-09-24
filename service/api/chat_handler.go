package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"new-wasa/service"
	"new-wasa/service/models"
)

// ChatHandler handles HTTP requests for chat operations
type ChatHandler struct {
	service *service.ChatService
}

// getPathParam extracts the last path parameter from the URL
func getPathParam(r *http.Request) string {
	return r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
}

// NewChatHandler creates a new chat handler
func NewChatHandler(service *service.ChatService) *ChatHandler {
	return &ChatHandler{
		service: service,
	}
}

// DoLogin handles user login
func (h *ChatHandler) DoLogin(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	response, err := h.service.Login(req.Name)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// SetMyUserName handles username update
func (h *ChatHandler) SetMyUserName(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := h.service.SetMyUserName(userID, req.Username); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Username updated successfully"})
}

// GetMyConversations handles getting user conversations
func (h *ChatHandler) GetMyConversations(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conversations, err := h.service.GetMyConversations(userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(conversations)
}

// GetConversation handles getting a specific conversation
func (h *ChatHandler) GetConversation(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conversationID := getPathParam(r)
	if conversationID == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	conversation, err := h.service.GetConversation(conversationID, userID)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(conversation)
}

// CreateConversation handles creating a new conversation
func (h *ChatHandler) CreateConversation(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	convType := req.Type
	if convType == "" {
		convType = "group"
	}

	conversation, err := h.service.CreateConversation(req.Name, convType, userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(conversation)
}

// SendMessage handles sending a message
func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		ConversationID string `json:"conversationId"`
		Body           string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	message, err := h.service.SendMessage(req.ConversationID, userID, req.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

// ForwardMessage handles forwarding a message
func (h *ChatHandler) ForwardMessage(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	messageID := getPathParam(r)
	if messageID == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var req struct {
		ConversationID string `json:"conversationId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	message, err := h.service.ForwardMessage(messageID, req.ConversationID, userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message)
}

// CommentMessage handles commenting on a message
func (h *ChatHandler) CommentMessage(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	messageID := getPathParam(r)
	if messageID == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var req struct {
		Comment string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	comment, err := h.service.CommentMessage(messageID, userID, req.Comment)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comment)
}

// UncommentMessage handles removing a comment from a message
func (h *ChatHandler) UncommentMessage(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	messageID := getPathParam(r)
	if messageID == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := h.service.UncommentMessage(messageID, userID); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment removed successfully"})
}

// DeleteMessage handles deleting a message
func (h *ChatHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	messageID := getPathParam(r)
	if messageID == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteMessage(messageID); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddToGroup handles adding a user to a group
func (h *ChatHandler) AddToGroup(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupID := getPathParam(r)
	if groupID == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var req struct {
		UserID string `json:"userId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := h.service.AddToGroup(groupID, req.UserID); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User added to group successfully"})
}

// LeaveGroup handles leaving a group
func (h *ChatHandler) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupID := getPathParam(r)
	if groupID == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := h.service.LeaveGroup(groupID, userID); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Left group successfully"})
}

// SetGroupName handles updating a group's name
func (h *ChatHandler) SetGroupName(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupID := getPathParam(r)
	if groupID == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := h.service.SetGroupName(groupID, req.Name, userID); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Group name updated successfully"})
}

// SetMyPhoto handles updating user's profile photo
func (h *ChatHandler) SetMyPhoto(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	photoData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := h.service.SetMyPhoto(userID, photoData); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":   "Photo updated successfully",
		"photoUrl": fmt.Sprintf("/photos/user_%s.jpg", userID),
	})
}

// SetGroupPhoto handles updating a group's photo
func (h *ChatHandler) SetGroupPhoto(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupID := getPathParam(r)
	if groupID == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	photoData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := h.service.SetGroupPhoto(groupID, userID, photoData); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":   "Group photo updated successfully",
		"photoUrl": fmt.Sprintf("/photos/group_%s.jpg", groupID),
	})
}
