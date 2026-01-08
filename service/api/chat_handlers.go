package api

import (
	"encoding/json"
	"net/http"
	"new-wasa/service/api/reqcontext"
	"new-wasa/service/database"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// listChats returns the list of users the requester follows (potential chat peers)
func (rt *_router) listChats(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	requester := extractBearer(r.Header.Get("Authorization"))
	if status := validateRequestingUser(ps.ByName("id"), requester); status != 0 {
		w.WriteHeader(status)
		return
	}
	convs, err := rt.db.ListConversations(database.User{IdUser: requester})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Wrap in an object to avoid top-level array responses (OpenAPI lint requirement).
	type conversationsResponse struct {
		Conversations []database.Conversation `json:"conversations"`
	}
	if err := json.NewEncoder(w).Encode(conversationsResponse{Conversations: convs}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// listMessages returns messages between requester and peer
func (rt *_router) listMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	requester := extractBearer(r.Header.Get("Authorization"))
	if status := validateRequestingUser(ps.ByName("id"), requester); status != 0 {
		w.WriteHeader(status)
		return
	}
	peer := ps.ByName("peer")
	var msgs []database.Message
	if strings.HasPrefix(peer, "g-") {
		groupID, err := strconv.ParseInt(strings.TrimPrefix(peer, "g-"), 10, 64)
		if err != nil || groupID <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		inGroup, err := rt.db.IsUserInGroup(groupID, database.User{IdUser: requester})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !inGroup {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		// Mark as read all group messages for this user.
		_ = rt.db.MarkGroupConversationRead(groupID, database.User{IdUser: requester})

		gmsgs, err := rt.db.ListGroupMessages(groupID, 100, 0)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, gm := range gmsgs {
			msg := database.Message{
				Id:       gm.Id,
				Sender:   gm.Sender,
				Receiver: peer,
				Body:     gm.Body,
				Date:     gm.Date,
			}
			if reactions, err := rt.db.ListGroupMessageReactions(groupID, msg.Id); err == nil {
				msg.Reactions = reactions
			}
			if msg.Sender == requester {
				if status, err := rt.db.GetGroupMessageCheckmarks(groupID, msg.Id); err == nil {
					msg.Status = status
				}
			}
			msgs = append(msgs, msg)
		}
	} else {
		// Mark as read all messages sent by peer to requester.
		_ = rt.db.MarkDirectConversationRead(database.User{IdUser: requester}, database.User{IdUser: peer})

		directMsgs, err := rt.db.ListMessages(database.User{IdUser: requester}, database.User{IdUser: peer}, 100, 0)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for i := range directMsgs {
			if reactions, err := rt.db.ListDirectMessageReactions(directMsgs[i].Id); err == nil {
				directMsgs[i].Reactions = reactions
			}
			if directMsgs[i].Sender == requester {
				if status, err := rt.db.GetDirectMessageCheckmarks(directMsgs[i].Id); err == nil {
					directMsgs[i].Status = status
				}
			}
		}
		msgs = directMsgs
	}

	// Wrap in an object to avoid top-level array responses (OpenAPI lint requirement).
	type messagesResponse struct {
		Messages []database.Message `json:"messages"`
	}

	if err := json.NewEncoder(w).Encode(messagesResponse{Messages: msgs}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// sendMessage sends a message from requester to peer
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	requester := extractBearer(r.Header.Get("Authorization"))
	if status := validateRequestingUser(ps.ByName("id"), requester); status != 0 {
		w.WriteHeader(status)
		return
	}
	peer := ps.ByName("peer")
	var body struct {
		Body string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.Body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	peerIsGroup := strings.HasPrefix(peer, "g-")
	if peerIsGroup {
		groupID, err := strconv.ParseInt(strings.TrimPrefix(peer, "g-"), 10, 64)
		if err != nil || groupID <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		inGroup, err := rt.db.IsUserInGroup(groupID, database.User{IdUser: requester})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !inGroup {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		_, err = rt.db.CreateGroupMessage(groupID, database.User{IdUser: requester}, body.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		// basic ban check: if requester is banned by peer, disallow
		banned, err := rt.db.BannedUserCheck(database.User{IdUser: requester}, database.User{IdUser: peer})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if banned {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		_, err = rt.db.CreateMessage(database.User{IdUser: requester}, database.User{IdUser: peer}, body.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
