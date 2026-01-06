package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"new-wasa/service/api/reqcontext"
	"new-wasa/service/database"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func parseGroupPeer(peer string) (int64, bool) {
	if !strings.HasPrefix(peer, "g-") {
		return 0, false
	}
	id, err := strconv.ParseInt(strings.TrimPrefix(peer, "g-"), 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	requester := extractBearer(r.Header.Get("Authorization"))
	if status := validateRequestingUser(ps.ByName("id"), requester); status != 0 {
		w.WriteHeader(status)
		return
	}

	messageID, err := strconv.ParseInt(ps.ByName("message_id"), 10, 64)
	if err != nil || messageID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	peer := ps.ByName("peer")
	if groupID, ok := parseGroupPeer(peer); ok {
		inGroup, err := rt.db.IsUserInGroup(groupID, database.User{IdUser: requester})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !inGroup {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Validate message exists in group
		if _, err := rt.db.GetGroupMessageInGroup(groupID, messageID); err != nil {
			if errors.Is(err, database.ErrMessageNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = rt.db.DeleteGroupMessage(groupID, messageID, database.User{IdUser: requester})
		if err != nil {
			if errors.Is(err, database.ErrMessageNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if errors.Is(err, database.ErrForbiddenMessageAction) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			ctx.Logger.WithError(err).Error("deleteMessage: db.DeleteGroupMessage error")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// direct chat message
	if _, err := rt.db.GetDirectMessageInConversation(database.User{IdUser: requester}, database.User{IdUser: peer}, messageID); err != nil {
		if errors.Is(err, database.ErrMessageNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = rt.db.DeleteDirectMessage(messageID, database.User{IdUser: requester})
	if err != nil {
		if errors.Is(err, database.ErrMessageNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, database.ErrForbiddenMessageAction) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		ctx.Logger.WithError(err).Error("deleteMessage: db.DeleteDirectMessage error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	requester := extractBearer(r.Header.Get("Authorization"))
	if status := validateRequestingUser(ps.ByName("id"), requester); status != 0 {
		w.WriteHeader(status)
		return
	}

	messageID, err := strconv.ParseInt(ps.ByName("message_id"), 10, 64)
	if err != nil || messageID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	type reqBody struct {
		Reaction string `json:"reaction"`
	}
	var rb reqBody
	if err := json.NewDecoder(r.Body).Decode(&rb); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rb.Reaction = strings.TrimSpace(rb.Reaction)
	if rb.Reaction == "" || len(rb.Reaction) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	peer := ps.ByName("peer")
	if groupID, ok := parseGroupPeer(peer); ok {
		inGroup, err := rt.db.IsUserInGroup(groupID, database.User{IdUser: requester})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !inGroup {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if _, err := rt.db.GetGroupMessageInGroup(groupID, messageID); err != nil {
			if errors.Is(err, database.ErrMessageNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := rt.db.SetGroupMessageReaction(groupID, messageID, database.User{IdUser: requester}, rb.Reaction); err != nil {
			ctx.Logger.WithError(err).Error("commentMessage: db.SetGroupMessageReaction error")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if _, err := rt.db.GetDirectMessageInConversation(database.User{IdUser: requester}, database.User{IdUser: peer}, messageID); err != nil {
		if errors.Is(err, database.ErrMessageNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := rt.db.SetDirectMessageReaction(messageID, database.User{IdUser: requester}, rb.Reaction); err != nil {
		ctx.Logger.WithError(err).Error("commentMessage: db.SetDirectMessageReaction error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	requester := extractBearer(r.Header.Get("Authorization"))
	if status := validateRequestingUser(ps.ByName("id"), requester); status != 0 {
		w.WriteHeader(status)
		return
	}

	messageID, err := strconv.ParseInt(ps.ByName("message_id"), 10, 64)
	if err != nil || messageID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	peer := ps.ByName("peer")
	if groupID, ok := parseGroupPeer(peer); ok {
		inGroup, err := rt.db.IsUserInGroup(groupID, database.User{IdUser: requester})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !inGroup {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if _, err := rt.db.GetGroupMessageInGroup(groupID, messageID); err != nil {
			if errors.Is(err, database.ErrMessageNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := rt.db.RemoveGroupMessageReaction(groupID, messageID, database.User{IdUser: requester}); err != nil {
			ctx.Logger.WithError(err).Error("uncommentMessage: db.RemoveGroupMessageReaction error")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if _, err := rt.db.GetDirectMessageInConversation(database.User{IdUser: requester}, database.User{IdUser: peer}, messageID); err != nil {
		if errors.Is(err, database.ErrMessageNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := rt.db.RemoveDirectMessageReaction(messageID, database.User{IdUser: requester}); err != nil {
		ctx.Logger.WithError(err).Error("uncommentMessage: db.RemoveDirectMessageReaction error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	requester := extractBearer(r.Header.Get("Authorization"))
	if status := validateRequestingUser(ps.ByName("id"), requester); status != 0 {
		w.WriteHeader(status)
		return
	}

	messageID, err := strconv.ParseInt(ps.ByName("message_id"), 10, 64)
	if err != nil || messageID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	type forwardReq struct {
		To string `json:"to"`
	}
	var fr forwardReq
	if err := json.NewDecoder(r.Body).Decode(&fr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fr.To = strings.TrimSpace(fr.To)
	if fr.To == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	peer := ps.ByName("peer")
	var body string

	// Read source message body
	if groupID, ok := parseGroupPeer(peer); ok {
		inGroup, err := rt.db.IsUserInGroup(groupID, database.User{IdUser: requester})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !inGroup {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		gm, err := rt.db.GetGroupMessageInGroup(groupID, messageID)
		if err != nil {
			if errors.Is(err, database.ErrMessageNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		body = gm.Body
	} else {
		m, err := rt.db.GetDirectMessageInConversation(database.User{IdUser: requester}, database.User{IdUser: peer}, messageID)
		if err != nil {
			if errors.Is(err, database.ErrMessageNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		body = m.Body
	}

	// Send to destination conversation
	if toGroupID, ok := parseGroupPeer(fr.To); ok {
		inGroup, err := rt.db.IsUserInGroup(toGroupID, database.User{IdUser: requester})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !inGroup {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if _, err := rt.db.CreateGroupMessage(toGroupID, database.User{IdUser: requester}, body); err != nil {
			ctx.Logger.WithError(err).Error("forwardMessage: db.CreateGroupMessage error")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// direct destination
	if !validIdentifier(fr.To) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	banned, err := rt.db.BannedUserCheck(database.User{IdUser: requester}, database.User{IdUser: fr.To})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if banned {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if _, err := rt.db.CreateMessage(database.User{IdUser: requester}, database.User{IdUser: fr.To}, body); err != nil {
		ctx.Logger.WithError(err).Error("forwardMessage: db.CreateMessage error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}


