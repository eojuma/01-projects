package routes

import (
	"context"
	"database/sql"
	"net/http"

	"social-network/backend/pkg/handlers"
	"social-network/backend/pkg/websocket"
)

// SetupRoutes registers application routes on the provided mux.
func SetupRoutes(mux *http.ServeMux, db *sql.DB, authHandler *handlers.AuthHandler, groupHandler *handlers.GroupHandler, profileHandler http.HandlerFunc, followHandler *handlers.FollowHandler, postHandler *handlers.PostHandler) {
	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/logout", authHandler.Logout)
	mux.HandleFunc("/api/auth/me", authHandler.GetCurrentUser)

	// Wire group handlers (simple dispatch by method)
	mux.HandleFunc("/api/groups", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			groupHandler.GetAllGroups(w, r)
			return
		}
		if r.Method == http.MethodPost {
			groupHandler.CreateGroup(w, r)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("GET /api/groups/{groupID}", groupHandler.GetGroupDetails)
	mux.HandleFunc("GET /api/groups/{groupID}/requests", groupHandler.GetJoinRequests)
	mux.HandleFunc("PATCH /api/groups/{groupID}", groupHandler.UpdateGroup)
	mux.HandleFunc("GET /api/groups/{groupID}/members", groupHandler.GetGroupMembersList)

	// Group invites
	mux.HandleFunc("POST /api/groups/{groupID}/invite", groupHandler.InviteUser)
	mux.HandleFunc("POST /api/groups/{groupID}/invites/{inviteID}/accept", groupHandler.AcceptInvite)
	mux.HandleFunc("POST /api/groups/{groupID}/invites/{inviteID}/decline", groupHandler.DeclineInvite)

	// Group join requests
	mux.HandleFunc("POST /api/groups/{groupID}/request", groupHandler.RequestToJoin)
	mux.HandleFunc("POST /api/groups/{groupID}/requests/{requestID}/approve", groupHandler.ApproveRequest)
	mux.HandleFunc("POST /api/groups/{groupID}/requests/{requestID}/decline", groupHandler.DeclineRequest)

	// Group posts
	mux.HandleFunc("POST /api/groups/{groupID}/posts", groupHandler.CreateGroupPost)

	// Group post feed
	mux.HandleFunc("GET /api/groups/{groupID}/posts", groupHandler.GetGroupFeed)
	// Group comments
	mux.HandleFunc("POST /api/groups/{groupID}/posts/{postID}/comments", groupHandler.CreateGroupComment)
	mux.HandleFunc("GET /api/groups/{groupID}/posts/{postID}/comments", groupHandler.GetGroupComments)
	// Group Events
	mux.HandleFunc("POST /api/groups/{groupID}/events", groupHandler.CreateEvent)
	mux.HandleFunc("GET /api/groups/{groupID}/events", groupHandler.GetEvents)
	mux.HandleFunc("POST /api/events/{eventID}/respond", groupHandler.RespondToEvent)
	mux.HandleFunc("GET /api/groups/{groupID}/messages", handlers.GetGroupMessagesHandler(db))

	mux.HandleFunc("GET /api/users", handlers.SearchUsersHandler(db))
	mux.HandleFunc("GET /api/users/{username}", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "mux.vars", map[string]string{
			"username": r.PathValue("username"),
		})
		profileHandler.ServeHTTP(w, r.WithContext(ctx))
	})

	// Profile update + privacy toggle (authenticated user only)
	mux.HandleFunc("PATCH /api/users/{id}/profile", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "mux.vars", map[string]string{
			"id": r.PathValue("id"),
		})
		handlers.UpdateProfileHandler(db).ServeHTTP(w, r.WithContext(ctx))
	})
	mux.HandleFunc("PATCH /api/users/{id}/privacy", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "mux.vars", map[string]string{
			"id": r.PathValue("id"),
		})
		handlers.UpdatePrivacyHandler(db).ServeHTTP(w, r.WithContext(ctx))
	})

	// websocket endpoint
	mux.HandleFunc("/ws", websocket.ServeWebSocket)

	// Follow endpoints
	mux.HandleFunc("POST /api/users/{id}/follow", followHandler.FollowUser)
	mux.HandleFunc("POST /api/follow-requests/{id}/respond", followHandler.RespondToFollowRequest)
	mux.HandleFunc("DELETE /api/users/{id}/unfollow", followHandler.Unfollow)
	mux.HandleFunc("GET /api/users/{id}/followers", followHandler.GetFollowers)
	mux.HandleFunc("GET /api/users/{id}/following", followHandler.GetFollowing)

	// Post endpoints
	mux.HandleFunc("POST /api/posts", postHandler.CreatePost)
	mux.HandleFunc("GET /api/posts/feed", postHandler.GetFeed)
	mux.HandleFunc("POST /api/posts/{id}/comments", postHandler.CreateComment)
	mux.HandleFunc("GET /api/posts/{id}/comments", postHandler.GetComments)
	mux.HandleFunc("GET /api/users/{id}/posts", postHandler.GetUserPosts)
	mux.HandleFunc("POST /api/posts/{id}/like", postHandler.ToggleLike)

	// Chat endpoints
	mux.HandleFunc("GET /api/conversations", handlers.GetConversationsHandler(db))
	mux.HandleFunc("GET /api/conversations/{id}/messages", handlers.GetMessagesHandler(db))
	mux.HandleFunc("GET /api/chat/user/{userID}", handlers.GetOrCreateConversationMessagesHandler(db))
	mux.HandleFunc("POST /api/chat/upload", handlers.ChatUploadHandler(db))

	// Notification endpoints
	mux.HandleFunc("GET /api/notifications", handlers.GetNotificationsHandler(db))
	mux.HandleFunc("POST /api/notifications/read", handlers.MarkNotificationsReadHandler(db))

	// Serve uploaded files
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))
}
