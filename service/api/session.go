package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"new-wasa/service/api/reqcontext"
	"new-wasa/service/database"
	"os"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) sessionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	type loginRequest struct {
		Name string `json:"name"`
	}
	type loginResponse struct {
		Identifier string `json:"identifier"`
	}

	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if !validIdentifier(req.Name) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// If a user with this nickname already exists, log them in and return its identifier.
	existingUser, found, err := rt.db.FindUserByNickname(req.Name)
	if err != nil {
		ctx.Logger.WithError(err).Error("session: error searching user by nickname")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if found {
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(loginResponse{Identifier: existingUser.IdUser})
		return
	}

	// Otherwise create a new user with a fresh identifier and set nickname=name.
	newIdentifier, err := generateUniqueIdentifier(rt.db)
	if err != nil {
		ctx.Logger.WithError(err).Error("session: error generating new identifier")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = rt.db.CreateUser(User{IdUser: newIdentifier}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("session: error creating user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = rt.db.ModifyNickname(User{IdUser: newIdentifier}.ToDatabase(), Nickname{Nickname: req.Name}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("session: error setting nickname")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create user's directories locally
	err = createUserFolder(newIdentifier, ctx)
	if err != nil {
		ctx.Logger.WithError(err).Error("session: can't create user's photo folder")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(loginResponse{Identifier: newIdentifier})
}

func generateUniqueIdentifier(db interface {
	CheckUser(a database.User) (bool, error)
}) (string, error) {
	for i := 0; i < 10; i++ {
		id, err := randomIdentifier16()
		if err != nil {
			return "", err
		}
		exists, err := db.CheckUser(User{IdUser: id}.ToDatabase())
		if err != nil {
			return "", err
		}
		if !exists {
			return id, nil
		}
	}
	return "", errors.New("could not generate a unique identifier")
}

func randomIdentifier16() (string, error) {
	b := make([]byte, 8) // 8 bytes -> 16 hex chars (fits VARCHAR(16))
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Function that creates a new subdir for the specified user
func createUserFolder(identifier string, ctx reqcontext.RequestContext) error {

	// Create the path media/useridentifier/ inside the project dir
	path := filepath.Join(photoFolder, identifier)

	// To the previously created path add the "photos" subdir
	err := os.MkdirAll(filepath.Join(path, "photos"), os.ModePerm)
	if err != nil {
		ctx.Logger.WithError(err).Error("session/createUserFolder:: error creating directories for user")
		return err
	}
	return nil
}
