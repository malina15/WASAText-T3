package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// Login enpoint
	rt.router.POST("/session", rt.wrap(rt.sessionHandler))

	// Search endpoint
	rt.router.GET("/users", rt.wrap(rt.getUsersQuery))

	// User Endpoint
	rt.router.PUT("/users/:id", rt.wrap(rt.putNickname))
	rt.router.GET("/users/:id", rt.wrap(rt.getUserProfile))
	rt.router.PUT("/users/:id/photo", rt.wrap(rt.setMyPhoto))
	rt.router.GET("/users/:id/photo", rt.wrap(rt.getUserPhoto))

	// Ban endpoint
	rt.router.PUT("/users/:id/banned_users/:banned_id", rt.wrap(rt.putBan))
	rt.router.DELETE("/users/:id/banned_users/:banned_id", rt.wrap(rt.deleteBan))

	// Followers endpoint
	rt.router.PUT("/users/:id/followers/:follower_id", rt.wrap(rt.putFollow))
	rt.router.DELETE("/users/:id/followers/:follower_id", rt.wrap(rt.deleteFollow))

	// Stream endpoint
	rt.router.GET("/users/:id/home", rt.wrap(rt.getHome))

	// Chat endpoints
	rt.router.GET("/users/:id/chats", rt.wrap(rt.listChats))
	rt.router.GET("/users/:id/chats/:peer/messages", rt.wrap(rt.listMessages))
	rt.router.POST("/users/:id/chats/:peer/messages", rt.wrap(rt.sendMessage))
	rt.router.DELETE("/users/:id/chats/:peer/messages/:message_id", rt.wrap(rt.deleteMessage))
	rt.router.POST("/users/:id/chats/:peer/messages/:message_id/comments", rt.wrap(rt.commentMessage))
	rt.router.DELETE("/users/:id/chats/:peer/messages/:message_id/comments", rt.wrap(rt.uncommentMessage))
	rt.router.POST("/users/:id/chats/:peer/messages/:message_id/forward", rt.wrap(rt.forwardMessage))

	// Group endpoints
	rt.router.POST("/users/:id/groups", rt.wrap(rt.createGroup))
	rt.router.PUT("/groups/:group_id/members/:member_id", rt.wrap(rt.addToGroup))
	rt.router.DELETE("/groups/:group_id/members/:member_id", rt.wrap(rt.leaveGroup))
	rt.router.PUT("/groups/:group_id", rt.wrap(rt.setGroupName))
	rt.router.PUT("/groups/:group_id/photo", rt.wrap(rt.setGroupPhoto))
	rt.router.GET("/groups/:group_id/photo", rt.wrap(rt.getGroupPhoto))

	// Photo Endpoint
	rt.router.POST("/users/:id/photos", rt.wrap(rt.postPhoto))
	rt.router.DELETE("/users/:id/photos/:photo_id", rt.wrap(rt.deletePhoto))
	rt.router.GET("/users/:id/photos/:photo_id", rt.wrap(rt.getPhoto))

	// Comments endpoint
	rt.router.POST("/users/:id/photos/:photo_id/comments", rt.wrap(rt.postComment))
	rt.router.DELETE("/users/:id/photos/:photo_id/comments/:comment_id", rt.wrap(rt.deleteComment))

	// Likes endpoint
	rt.router.PUT("/users/:id/photos/:photo_id/likes/:like_id", rt.wrap(rt.putLike))
	rt.router.DELETE("/users/:id/photos/:photo_id/likes/:like_id", rt.wrap(rt.deleteLike))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
