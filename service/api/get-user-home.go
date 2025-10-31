package api

import (
	"encoding/json"
	"net/http"
	"new-wasa/service/api/reqcontext"
	"new-wasa/service/database"

	"github.com/julienschmidt/httprouter"
)

// This function retrieves all the photos of the people that the user is following
func (rt *_router) getHome(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")
	identifier := extractBearer(r.Header.Get("Authorization"))

	// A user can only see his/her home
	valid := validateRequestingUser(ps.ByName("id"), identifier)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	followers, err := rt.db.GetFollowing(User{IdUser: identifier}.ToDatabase())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	photos := make([]database.Photo, 0)
	for _, follower := range followers {

		followerPhoto, err := rt.db.GetPhotosList(
			User{IdUser: identifier}.ToDatabase(),
			User{IdUser: follower.IdUser}.ToDatabase())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// If user doesn't follow anyone or stream is empty, show recent posts from all users
		if len(photos) == 0 {
			allUsers, err := rt.db.GetAllUsers()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			for _, u := range allUsers {
				// skip banned relationships
				banned, berr := rt.db.BannedUserCheck(User{IdUser: identifier}.ToDatabase(), User{IdUser: u.IdUser}.ToDatabase())
				if berr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if banned {
					continue
				}
				lst, lerr := rt.db.GetPhotosList(User{IdUser: identifier}.ToDatabase(), User{IdUser: u.IdUser}.ToDatabase())
				if lerr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				// limit per user
				for i, p := range lst {
					if i >= database.PhotosPerUserHome {
						break
					}
					photos = append(photos, p)
				}
			}
		}

		for i, photo := range followerPhoto {
			if i >= database.PhotosPerUserHome {
				break
			}
			photos = append(photos, photo)
		}

	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(photos); err != nil {
		ctx.Logger.WithError(err).Error("getHome/Encode: failed to encode photos json")
		return
	}
}
