package api

import (
    "encoding/json"
    "net/http"
    "new-wasa/service/api/reqcontext"
    "new-wasa/service/database"

    "github.com/julienschmidt/httprouter"
)

// listChats returns the list of users the requester follows (potential chat peers)
func (rt *_router) listChats(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    w.Header().Set("Content-Type", "application/json")
    requester := extractBearer(r.Header.Get("Authorization"))
    if validateRequestingUser(ps.ByName("id"), requester) != 0 {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    following, err := rt.db.GetFollowing(database.User{IdUser: requester})
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    if err := json.NewEncoder(w).Encode(following); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
    }
}

// listMessages returns messages between requester and peer
func (rt *_router) listMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    w.Header().Set("Content-Type", "application/json")
    requester := extractBearer(r.Header.Get("Authorization"))
    if validateRequestingUser(ps.ByName("id"), requester) != 0 {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    peer := ps.ByName("peer")
    msgs, err := rt.db.ListMessages(database.User{IdUser: requester}, database.User{IdUser: peer}, 100, 0)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    if err := json.NewEncoder(w).Encode(msgs); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
    }
}

// sendMessage sends a message from requester to peer
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    w.Header().Set("Content-Type", "application/json")
    requester := extractBearer(r.Header.Get("Authorization"))
    if validateRequestingUser(ps.ByName("id"), requester) != 0 {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    peer := ps.ByName("peer")
    var body struct{ Body string `json:"body"` }
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.Body) == 0 {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
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
    w.WriteHeader(http.StatusNoContent)
}


