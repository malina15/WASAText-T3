package api

import (
	"errors"
	"io"
	"net/http"
	"new-wasa/service/api/reqcontext"
	"new-wasa/service/database"
	"os"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// setMyPhoto uploads/updates the user's profile photo.
func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	requester := extractBearer(r.Header.Get("Authorization"))
	if validateRequestingUser(ps.ByName("id"), requester) != 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("setMyPhoto: error reading body")
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

	userDir := filepath.Join(photoFolder, requester)
	if err := os.MkdirAll(userDir, os.ModePerm); err != nil {
		ctx.Logger.WithError(err).Error("setMyPhoto: error creating user dir")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	filename := "avatar." + ext
	fullPath := filepath.Join(userDir, filename)
	if err := os.WriteFile(fullPath, data, 0o644); err != nil {
		ctx.Logger.WithError(err).Error("setMyPhoto: error writing file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	relPath := filepath.Join(requester, filename)
	err = rt.db.SetUserPhotoPath(database.User{IdUser: requester}, relPath)
	if err != nil {
		ctx.Logger.WithError(err).Error("setMyPhoto: db.SetUserPhotoPath error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// getUserPhoto serves a user's profile photo (requires requester to be logged in and not banned by target user).
func (rt *_router) getUserPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	requester := extractBearer(r.Header.Get("Authorization"))
	if isNotLogged(requester) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	target := ps.ByName("id")
	if !validIdentifier(target) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// If requester is banned by target, disallow access.
	banned, err := rt.db.BannedUserCheck(database.User{IdUser: requester}, database.User{IdUser: target})
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserPhoto: db.BannedUserCheck error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if banned {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	photoPath, err := rt.db.GetUserPhotoPath(database.User{IdUser: target})
	if err != nil {
		if errors.Is(err, database.ErrUserPhotoNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("getUserPhoto: db.GetUserPhotoPath error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if strings.TrimSpace(photoPath) == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, filepath.Join(photoFolder, photoPath))
}


