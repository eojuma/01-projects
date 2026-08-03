package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"

	"social-network/backend/pkg/repositories"
	"social-network/backend/pkg/websocket"
)

type FollowHandler struct {
	Repo *repositories.Repository
	Hub  *websocket.Hub
}

func NewFollowHandler(repo *repositories.Repository, hub *websocket.Hub) *FollowHandler {
	return &FollowHandler{Repo: repo, Hub: hub}
}

func (h *FollowHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	requesterID, ok := r.Context().Value("userID").(string)
	if !ok || requesterID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	targetID := r.PathValue("id")
	if targetID == "" {
		http.Error(w, "missing target user id", http.StatusBadRequest)
		return
	}

	requestID := uuid.Must(uuid.NewV4()).String()
	notificationID := uuid.Must(uuid.NewV4()).String()

	result, _, err := h.Repo.HandleFollowRequest(requestID, notificationID, requesterID, targetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *FollowHandler) RespondToFollowRequest(w http.ResponseWriter, r *http.Request) {
	targetID, ok := r.Context().Value("userID").(string)
	if !ok || targetID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	requestID := r.PathValue("id")
	if requestID == "" {
		http.Error(w, "missing request id", http.StatusBadRequest)
		return
	}

	action := r.URL.Query().Get("action")
	if action != "accept" && action != "decline" {
		http.Error(w, "action must be accept or decline", http.StatusBadRequest)
		return
	}

	result, err := h.Repo.RespondToFollowRequest(requestID, targetID, action)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *FollowHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	followerID, ok := r.Context().Value("userID").(string)
	if !ok || followerID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	followingID := r.PathValue("id")
	if followingID == "" {
		http.Error(w, "missing target user id", http.StatusBadRequest)
		return
	}

	err := h.Repo.Unfollow(followerID, followingID)
	if err != nil {
		if err.Error() == "not found" {
			http.Error(w, "follow relationship not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *FollowHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	sessionUserID, ok := r.Context().Value("userID").(string)
	if !ok || sessionUserID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	targetID := r.PathValue("id")
	if targetID == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
	}

	page, perPage := parsePagination(r)

	result, err := h.Repo.GetFollowers(targetID, sessionUserID, page, perPage)
	if err != nil {
		if err.Error() == "profile is private" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if err.Error() == "user not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *FollowHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	sessionUserID, ok := r.Context().Value("userID").(string)
	if !ok || sessionUserID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	targetID := r.PathValue("id")
	if targetID == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
	}

	page, perPage := parsePagination(r)

	result, err := h.Repo.GetFollowing(targetID, sessionUserID, page, perPage)
	if err != nil {
		if err.Error() == "profile is private" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if err.Error() == "user not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func parsePagination(r *http.Request) (page, perPage int) {
	page = 1
	perPage = 20
	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if pp := r.URL.Query().Get("per_page"); pp != "" {
		if v, err := strconv.Atoi(pp); err == nil && v > 0 && v <= 100 {
			perPage = v
		}
	}
	return
}
