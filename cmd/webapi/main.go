package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"new-wasa/service"
	"new-wasa/service/api"
	"new-wasa/service/store"
)

func main() {
	// Load configuration
	config := LoadConfig()
	
	// Initialize SQLite store
	store, err := store.NewSQLiteChatStore(config.Database.Path)
	if err != nil {
		log.Fatal("Failed to initialize store:", err)
	}

	// Initialize service and handlers
	chatService := service.NewChatService(store)
	h := api.NewChatHandler(chatService)

	// Create router
	router := httprouter.New()

	// CORS middleware wrapper for httprouter
	corsWrapper := SimpleCORSHandler

	// Combined middleware for protected routes
	protectedWrapper := func(handler http.HandlerFunc) httprouter.Handle {
		return corsWrapper(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			token := authHeader[7:] // Remove "Bearer " prefix

			// Validate token through the store
			user, err := store.GetUserByToken(token)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Add user to request context
			r.Header.Set("X-User-ID", user.ID)
			handler.ServeHTTP(w, r)
		})
	}

	// Public routes (no auth required)
	router.POST("/session", corsWrapper(h.DoLogin))
	router.OPTIONS("/session", corsWrapper(func(w http.ResponseWriter, r *http.Request) {}))

	// Protected routes (auth required)
	router.PUT("/user/username", protectedWrapper(h.SetMyUserName))
	router.OPTIONS("/user/username", corsWrapper(func(w http.ResponseWriter, r *http.Request) {}))
	
	router.POST("/conversations", protectedWrapper(h.CreateConversation))
	router.GET("/conversations", protectedWrapper(h.GetMyConversations))
	router.OPTIONS("/conversations", corsWrapper(func(w http.ResponseWriter, r *http.Request) {}))
	
	router.GET("/conversations/:id", protectedWrapper(h.GetConversation))
	router.OPTIONS("/conversations/:id", corsWrapper(func(w http.ResponseWriter, r *http.Request) {}))
	
	router.POST("/messages", protectedWrapper(h.SendMessage))
	router.OPTIONS("/messages", corsWrapper(func(w http.ResponseWriter, r *http.Request) {}))
	
	router.POST("/messages/:id/forward", protectedWrapper(h.ForwardMessage))
	router.OPTIONS("/messages/:id/forward", corsWrapper(func(w http.ResponseWriter, r *http.Request) {}))
	
	router.POST("/messages/:id/comment", protectedWrapper(h.CommentMessage))
	router.OPTIONS("/messages/:id/comment", corsWrapper(func(w http.ResponseWriter, r *http.Request) {}))
	
	router.POST("/messages/:id/uncomment", protectedWrapper(h.UncommentMessage))
	router.OPTIONS("/messages/:id/uncomment", corsWrapper(func(w http.ResponseWriter, r *http.Request) {}))
	
	router.DELETE("/messages/:id", protectedWrapper(h.DeleteMessage))
	router.OPTIONS("/messages/:id", corsWrapper(func(w http.ResponseWriter, r *http.Request) {}))
	
	router.POST("/groups/:id/add", protectedWrapper(h.AddToGroup))
	router.OPTIONS("/groups/:id/add", corsWrapper(func(w http.ResponseWriter, r *http.Request) {}))
	
	router.POST("/groups/:id/leave", protectedWrapper(h.LeaveGroup))
	router.OPTIONS("/groups/:id/leave", corsWrapper(func(w http.ResponseWriter, r *http.Request) {}))
	
	router.PUT("/groups/:id/name", protectedWrapper(h.SetGroupName))
	router.OPTIONS("/groups/:id/name", corsWrapper(func(w http.ResponseWriter, r *http.Request) {}))
	
	router.PUT("/user/photo", protectedWrapper(h.SetMyPhoto))
	router.OPTIONS("/user/photo", corsWrapper(func(w http.ResponseWriter, r *http.Request) {}))
	
	router.PUT("/groups/:id/photo", protectedWrapper(h.SetGroupPhoto))
	router.OPTIONS("/groups/:id/photo", corsWrapper(func(w http.ResponseWriter, r *http.Request) {}))

	// Register webui if built with webui tag
	registerWebUI()

	serverAddr := config.Server.Host + ":" + config.Server.Port
	log.Printf("Server listening on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}
