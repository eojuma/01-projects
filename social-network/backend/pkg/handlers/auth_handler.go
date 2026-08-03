package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"social-network/backend/pkg/auth"
	"social-network/backend/pkg/models"
	"social-network/backend/pkg/repositories"
	"social-network/backend/pkg/validation"
)

type AuthHandler struct {
	DB             *sql.DB
	SessionRepo    repositories.SessionRepository
	CookieName     string
	Timeout        time.Duration
	CookieDomain   string
	CookieSecure   bool
	CookieSameSite http.SameSite
}

type registerRequest struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
	Nickname    string `json:"nickname"`
	Bio         string `json:"bio"`
	AvatarURL   string `json:"avatar_url"`
}

type loginRequest struct {
	EmailOrUsername string `json:"email_or_username"`
	Password        string `json:"password"`
}

// GetCurrentUser returns the current authenticated user's info
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get userID from context (set by session middleware)
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := repositories.GetUserByIDString(h.DB, userID)
	if err != nil {
		if err == repositories.ErrUserNotFound {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Return user info without password
	resp := models.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Username:    user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Bio:         user.Bio,
		AvatarURL:   user.AvatarURL,
		Nickname:    user.Nickname,
		DateOfBirth: user.DateOfBirth,
		IsPrivate:   user.IsPrivate,
		CreatedAt:   user.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if !validation.ValidateEmail(req.Email) {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}
	if !validation.ValidateUsername(req.Username) {
		http.Error(w, "invalid username", http.StatusBadRequest)
		return
	}
	if !validation.ValidatePassword(req.Password) {
		http.Error(w, "password too short", http.StatusBadRequest)
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Bio:          req.Bio,
		AvatarURL:    req.AvatarURL,
		Nickname:     req.Nickname,
		DateOfBirth:  req.DateOfBirth,
	}

	if err := repositories.CreateUser(h.DB, user); err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "users.email") {
			http.Error(w, "email already registered", http.StatusBadRequest)
			return
		}
		if strings.Contains(errStr, "users.username") {
			http.Error(w, "username already taken", http.StatusBadRequest)
			return
		}
		http.Error(w, "could not create user", http.StatusInternalServerError)
		return
	}

	resp := models.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Username:    user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Bio:         user.Bio,
		AvatarURL:   user.AvatarURL,
		Nickname:    user.Nickname,
		DateOfBirth: user.DateOfBirth,
		IsPrivate:   user.IsPrivate,
		CreatedAt:   user.CreatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var user *models.User
	if validation.ValidateEmail(req.EmailOrUsername) {
		u, err := repositories.GetUserByEmail(h.DB, req.EmailOrUsername)
		if err == nil {
			user = u
		}
	}
	if user == nil {
		u, err := repositories.GetUserByUsername(h.DB, req.EmailOrUsername)
		if err == nil {
			user = u
		}
	}
	if user == nil {
		http.Error(w, "username or email not found", http.StatusUnauthorized)
		return
	}

	if err := auth.ComparePasswordHash(user.PasswordHash, req.Password); err != nil {
		http.Error(w, "incorrect password", http.StatusUnauthorized)
		return
	}

	sid, err := auth.GenerateSessionID()
	if err != nil {
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(h.Timeout)
	if err := h.SessionRepo.Create(sid, (user.ID), expiresAt); err != nil {
		http.Error(w, "failed to store session", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     h.CookieName,
		Value:    sid,
		Expires:  expiresAt,
		HttpOnly: true,
		Path:     "/",
		Domain:   h.CookieDomain,
		Secure:   h.CookieSecure,
		SameSite: h.CookieSameSite,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie(h.CookieName)
	if err == nil {
		_ = h.SessionRepo.Delete(cookie.Value)
		http.SetCookie(w, &http.Cookie{
			Name:     h.CookieName,
			Value:    "",
			Expires:  time.Unix(0, 0),
			Path:     "/",
			Domain:   h.CookieDomain,
			HttpOnly: true,
			Secure:   h.CookieSecure,
			SameSite: h.CookieSameSite,
		})
	}
	w.WriteHeader(http.StatusNoContent)
}
