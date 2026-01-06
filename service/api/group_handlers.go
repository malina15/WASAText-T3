package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"new-wasa/service/api/reqcontext"
	"new-wasa/service/database"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// createGroup creates a new group with the requesting user as member.
func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	requester := extractBearer(r.Header.Get("Authorization"))
	if validateRequestingUser(ps.ByName("id"), requester) != 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	type createGroupRequest struct {
		Name    string   `json:"name"`
		Members []string `json:"members"`
	}
	type createGroupResponse struct {
		GroupID int64 `json:"group_id"`
	}

	var req createGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if len(req.Name) == 0 || len(req.Name) > 32 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	members := make([]database.User, 0, len(req.Members))
	for _, raw := range req.Members {
		id := strings.TrimSpace(raw)
		if id == "" {
			continue
		}
		if !validIdentifier(id) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		exists, err := rt.db.CheckUser(database.User{IdUser: id})
		if err != nil {
			ctx.Logger.WithError(err).Error("createGroup: db.CheckUser error")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		members = append(members, database.User{IdUser: id})
	}

	groupID, err := rt.db.CreateGroup(database.User{IdUser: requester}, req.Name, members)
	if err != nil {
		ctx.Logger.WithError(err).Error("createGroup: db.CreateGroup error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(createGroupResponse{GroupID: groupID})
}

// addToGroup adds a user to a group (only group members can add).
func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	requester := extractBearer(r.Header.Get("Authorization"))
	if isNotLogged(requester) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	groupID, err := strconv.ParseInt(ps.ByName("group_id"), 10, 64)
	if err != nil || groupID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	memberID := strings.TrimSpace(ps.ByName("member_id"))
	if !validIdentifier(memberID) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inGroup, err := rt.db.IsUserInGroup(groupID, database.User{IdUser: requester})
	if err != nil {
		ctx.Logger.WithError(err).Error("addToGroup: db.IsUserInGroup error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !inGroup {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	exists, err := rt.db.CheckUser(database.User{IdUser: memberID})
	if err != nil {
		ctx.Logger.WithError(err).Error("addToGroup: db.CheckUser error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = rt.db.AddUserToGroup(groupID, database.User{IdUser: memberID})
	if err != nil {
		if errors.Is(err, database.ErrGroupNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, database.ErrUserAlreadyInGroup) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		ctx.Logger.WithError(err).Error("addToGroup: db.AddUserToGroup error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// leaveGroup removes the requesting user from the group.
func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	requester := extractBearer(r.Header.Get("Authorization"))
	if isNotLogged(requester) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	groupID, err := strconv.ParseInt(ps.ByName("group_id"), 10, 64)
	if err != nil || groupID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	memberID := strings.TrimSpace(ps.ByName("member_id"))
	if memberID != requester {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = rt.db.RemoveUserFromGroup(groupID, database.User{IdUser: requester})
	if err != nil {
		if errors.Is(err, database.ErrGroupNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, database.ErrUserNotInGroup) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		ctx.Logger.WithError(err).Error("leaveGroup: db.RemoveUserFromGroup error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// setGroupName updates the group's name (only group members can update).
func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	requester := extractBearer(r.Header.Get("Authorization"))
	if isNotLogged(requester) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	groupID, err := strconv.ParseInt(ps.ByName("group_id"), 10, 64)
	if err != nil || groupID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inGroup, err := rt.db.IsUserInGroup(groupID, database.User{IdUser: requester})
	if err != nil {
		ctx.Logger.WithError(err).Error("setGroupName: db.IsUserInGroup error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !inGroup {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	type setGroupNameRequest struct {
		Name string `json:"name"`
	}
	var req setGroupNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if len(req.Name) == 0 || len(req.Name) > 32 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.SetGroupName(groupID, req.Name)
	if err != nil {
		if errors.Is(err, database.ErrGroupNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("setGroupName: db.SetGroupName error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// setGroupPhoto sets the group photo (only group members can update).
func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	requester := extractBearer(r.Header.Get("Authorization"))
	if isNotLogged(requester) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	groupID, err := strconv.ParseInt(ps.ByName("group_id"), 10, 64)
	if err != nil || groupID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inGroup, err := rt.db.IsUserInGroup(groupID, database.User{IdUser: requester})
	if err != nil {
		ctx.Logger.WithError(err).Error("setGroupPhoto: db.IsUserInGroup error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !inGroup {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("setGroupPhoto: error reading body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	contentType := http.DetectContentType(data)
	ext := ""
	switch contentType {
	case "image/jpeg":
		ext = "jpg"
	case "image/png":
		ext = "png"
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	groupDir := filepath.Join(photoFolder, "groups", strconv.FormatInt(groupID, 10))
	if err := os.MkdirAll(groupDir, os.ModePerm); err != nil {
		ctx.Logger.WithError(err).Error("setGroupPhoto: error creating group dir")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	filename := "photo." + ext
	fullPath := filepath.Join(groupDir, filename)
	if err := os.WriteFile(fullPath, data, 0o644); err != nil {
		ctx.Logger.WithError(err).Error("setGroupPhoto: error writing file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	relPath := filepath.Join("groups", strconv.FormatInt(groupID, 10), filename)
	err = rt.db.SetGroupPhotoPath(groupID, relPath)
	if err != nil {
		if errors.Is(err, database.ErrGroupNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("setGroupPhoto: db.SetGroupPhotoPath error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// getGroupPhoto serves the group photo to group members.
func (rt *_router) getGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	requester := extractBearer(r.Header.Get("Authorization"))
	if isNotLogged(requester) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	groupID, err := strconv.ParseInt(ps.ByName("group_id"), 10, 64)
	if err != nil || groupID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inGroup, err := rt.db.IsUserInGroup(groupID, database.User{IdUser: requester})
	if err != nil {
		ctx.Logger.WithError(err).Error("getGroupPhoto: db.IsUserInGroup error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !inGroup {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	group, err := rt.db.GetGroup(groupID)
	if err != nil {
		if errors.Is(err, database.ErrGroupNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("getGroupPhoto: db.GetGroup error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if strings.TrimSpace(group.PhotoPath) == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filepath.Join(photoFolder, group.PhotoPath))
}


