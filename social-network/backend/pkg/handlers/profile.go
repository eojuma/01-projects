package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"social-network/backend/pkg/repositories"
)

type updateProfileRequest struct {
	Username  *string `json:"username"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Bio       *string `json:"bio"`
	AvatarURL *string `json:"avatar_url"`
	Nickname  *string `json:"nickname"`
	// Updating privacy uses a separate endpoint.
}

type updateProfileResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
	Nickname  string `json:"nickname"`
	IsPrivate bool   `json:"is_private"`
	UpdatedAt string `json:"updated_at"`
}

type updatePrivacyResponse struct {
	ID        string `json:"id"`
	IsPrivate bool   `json:"is_private"`
	UpdatedAt string `json:"updated_at"`
}

// UpdateProfileHandler handles PATCH /api/users/{id}/profile
// Only the authenticated user can update their own profile.
func UpdateProfileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		vars := muxVars(r)
		idStr, ok := vars["id"]
		if !ok || idStr == "" {
			http.Error(w, "missing user id", http.StatusBadRequest)
			return
		}
		targetID, err := strconv.Atoi(idStr)
		if err != nil || targetID <= 0 {
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}

		authID, ok := authUserIDFromContext(r)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if authID != strconv.Itoa(targetID) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		var req updateProfileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		// Basic validation: allow empty body but require at least one field.
		if req.Username == nil && req.FirstName == nil && req.LastName == nil && req.Bio == nil && req.AvatarURL == nil && req.Nickname == nil {
			http.Error(w, "no fields to update", http.StatusBadRequest)
			return
		}

		updated, err := repositories.UpdateUserProfile(db, targetID, repositories.UpdateUserProfileParams{
			Username:  req.Username,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Bio:       req.Bio,
			AvatarURL: req.AvatarURL,
			Nickname:  req.Nickname,
		})
		if err != nil {
			if errors.Is(err, repositories.ErrUserNotFound) {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}
			http.Error(w, "could not update profile", http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, updateProfileResponse{
			ID:        updated.ID,
			Username:  updated.Username,
			FirstName: updated.FirstName,
			LastName:  updated.LastName,
			Bio:       updated.Bio,
			AvatarURL: updated.AvatarURL,
			Nickname:  updated.Nickname,
			IsPrivate: updated.IsPrivate,
			UpdatedAt: updated.UpdatedAt.Format(time.RFC3339),
		})
	}
}

// UpdatePrivacyHandler handles PATCH /api/users/{id}/privacy
// Only the authenticated user can toggle their privacy.
func UpdatePrivacyHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		vars := muxVars(r)
		idStr, ok := vars["id"]
		if !ok || idStr == "" {
			http.Error(w, "missing user id", http.StatusBadRequest)
			return
		}
		targetID, err := strconv.Atoi(idStr)
		if err != nil || targetID <= 0 {
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}

		authID, ok := authUserIDFromContext(r)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if authID != strconv.Itoa(targetID) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		var req struct {
			IsPrivate bool `json:"is_private"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		updated, err := repositories.SetUserPrivacy(db, targetID, req.IsPrivate)
		if err != nil {
			if errors.Is(err, repositories.ErrUserNotFound) {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}
			http.Error(w, "could not update privacy", http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, updatePrivacyResponse{
			ID:        updated.ID,
			IsPrivate: updated.IsPrivate,
			UpdatedAt: updated.UpdatedAt.Format(time.RFC3339),
		})
	}
}
