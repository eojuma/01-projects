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

type publicProfileResponse struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	AvatarURL    string `json:"avatar_url,omitempty"`
	IsPrivate    bool   `json:"is_private"`
	CreatedAt    string `json:"created_at"`
	FollowStatus string `json:"follow_status,omitempty"`
}

type privateProfileResponse struct {
	publicProfileResponse
	Email     string `json:"email,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Bio       string `json:"bio,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

func GetProfileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := muxVars(r)
		username, ok := vars["username"]
		if !ok || username == "" {
			http.Error(w, "missing username", http.StatusBadRequest)
			return
		}

		user, err := repositories.GetUserByUsername(db, username)
		if err != nil {
			if errors.Is(err, repositories.ErrUserNotFound) {
				user, err = repositories.GetUserByIDString(db, username)
				if err != nil {
					if errors.Is(err, repositories.ErrUserNotFound) {
						http.Error(w, "user not found", http.StatusNotFound)
						return
					}
					http.Error(w, "internal error", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		}

		authID, _ := authUserIDFromContext(r)
		isFollower := false
		followStatus := ""
		if authID != "" && authID != user.ID {
			var count int
			err := db.QueryRow(
				`SELECT COUNT(*) FROM followers WHERE follower_id = ? AND following_id = ?`,
				authID, user.ID,
			).Scan(&count)
			if err == nil && count > 0 {
				isFollower = true
				followStatus = "following"
			} else {
				var pendingCount int
				_ = db.QueryRow(
					`SELECT COUNT(*) FROM follow_requests WHERE requester_id = ? AND target_id = ? AND status = 'pending'`,
					authID, user.ID,
				).Scan(&pendingCount)
				if pendingCount > 0 {
					followStatus = "pending"
				} else {
					followStatus = "none"
				}
			}
		}

		if user.IsPrivate && authID != user.ID && !isFollower {

			resp := publicProfileResponse{
				ID:           user.ID,
				Username:     user.Username,
				AvatarURL:    user.AvatarURL,
				IsPrivate:    user.IsPrivate,
				CreatedAt:    user.CreatedAt.Format(time.RFC3339),
				FollowStatus: followStatus,
			}
			writeJSON(w, http.StatusOK, resp)
			return
		}

		resp := privateProfileResponse{
			publicProfileResponse: publicProfileResponse{
				ID:           user.ID,
				Username:     user.Username,
				AvatarURL:    user.AvatarURL,
				IsPrivate:    user.IsPrivate,
				CreatedAt:    user.CreatedAt.Format(time.RFC3339),
				FollowStatus: followStatus,
			},
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Bio:       user.Bio,
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		}
		writeJSON(w, http.StatusOK, resp)
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func authUserIDFromContext(r *http.Request) (string, bool) {
	val := r.Context().Value("user_id")
	if val == nil {
		val = r.Context().Value("userID")
	}
	if val == nil {
		return "", false
	}
	switch v := val.(type) {
	case string:
		return v, true
	case int:
		return strconv.Itoa(v), true
	case int32:
		return strconv.FormatInt(int64(v), 10), true
	case int64:
		return strconv.FormatInt(v, 10), true
	case float64:
		return strconv.FormatInt(int64(v), 10), true
	default:
		return "", false
	}
}

func muxVars(r *http.Request) map[string]string {
	if v := r.Context().Value("mux.vars"); v != nil {
		if m, ok := v.(map[string]string); ok {
			return m
		}
	}
	return map[string]string{}
}

func SearchUsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		currentUserID, ok := authUserIDFromContext(r)
		if !ok || currentUserID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		query := r.URL.Query().Get("q")
		var users interface{}
		var err error
		if r.URL.Query().Get("chat") == "true" {
			users, err = repositories.SearchChatUsers(db, currentUserID, query)
		} else {
			users, err = repositories.SearchUsers(db, currentUserID, query)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, users)
	}
}
