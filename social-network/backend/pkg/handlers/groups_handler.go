package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"social-network/backend/pkg/models"
	"social-network/backend/pkg/repositories"
)

type GroupHandler struct {
	DB *sql.DB
}

func getSessionUser(r *http.Request) (string, error) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("unauthorized")
	}
	return userID, nil
}

func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getSessionUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	group, err := repositories.CreateGroup(h.DB, req, userID)
	if err != nil {
		http.Error(w, "could not create group", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(group)
}

func (h *GroupHandler) GetGroupDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getSessionUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	groupID := r.PathValue("groupID")
	if groupID == "" {
		groupID = r.PathValue("id")
	}
	if groupID == "" {
		http.Error(w, "missing group id", http.StatusBadRequest)
		return
	}

	group, err := repositories.GetGroupByID(h.DB, groupID)
	if err != nil {
		http.Error(w, "could not fetch group", http.StatusInternalServerError)
		return
	}
	if group == nil {
		http.Error(w, "group not found", http.StatusNotFound)
		return
	}

	isMember, err := repositories.IsMember(h.DB, groupID, userID)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	var isRequested bool
	err = h.DB.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM group_join_requests WHERE group_id = ? AND user_id = ? AND status = 'pending')`,
		groupID, userID,
	).Scan(&isRequested)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	isCreator := group.CreatorID == userID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":           group.ID,
		"title":        group.Title,
		"description":  group.Description,
		"creator_id":   group.CreatorID,
		"created_at":   group.CreatedAt,
		"is_member":    isMember,
		"is_creator":   isCreator,
		"is_requested": isRequested,
	})
}

func (h *GroupHandler) GetJoinRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getSessionUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	groupID := r.PathValue("groupID")
	if groupID == "" {
		http.Error(w, "missing group id", http.StatusBadRequest)
		return
	}

	// Verify current user is group creator
	isCreator, err := repositories.IsGroupCreator(h.DB, groupID, userID)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	if !isCreator {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	reqs, err := repositories.GetPendingJoinRequests(h.DB, groupID)
	if err != nil {
		http.Error(w, "could not fetch join requests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reqs)
}

func (h *GroupHandler) GetAllGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getSessionUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	pagestr := r.URL.Query().Get("page")
	limitstr := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")

	page := 1
	limit := 20

	if pagestr != "" {
		if p, err := strconv.Atoi(pagestr); err == nil && p > 0 {
			page = p
		}
	}

	if limitstr != "" {
		if y, err := strconv.Atoi(limitstr); err == nil && y > 0 {
			limit = y
		}
	}

	groups, err := repositories.GetAllGroups(h.DB, userID, page, limit, search)
	if err != nil {
		http.Error(w, "could not fetch groups", http.StatusInternalServerError)
		return
	}

	if groups == nil {
		groups = []models.GroupWithCount{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

func (h *GroupHandler) InviteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getSessionUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	groupID := r.PathValue("groupID")
	if groupID == "" {
		http.Error(w, "group id is required", http.StatusBadRequest)
		return
	}

	var body struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if body.UserID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	invite, err := repositories.CreateInvite(h.DB, groupID, userID, body.UserID)
	if err != nil {
		if err.Error() == "only group members can send invites" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if err.Error() == "user is already a member" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "could not create invite", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(invite)
}

func (h *GroupHandler) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getSessionUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	inviteID := r.PathValue("inviteID")
	if inviteID == "" {
		http.Error(w, "invite id is required", http.StatusBadRequest)
		return
	}

	if err := repositories.UpdateInviteStatus(h.DB, inviteID, userID, "accepted"); err != nil {
		if err.Error() == "unauthorized" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if err.Error() == "invite not found" {
			http.Error(w, "invite not found", http.StatusNotFound)
			return
		}
		http.Error(w, "could not accept invite", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "invite accepted"}`))
}

func (h *GroupHandler) DeclineInvite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getSessionUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	inviteID := r.PathValue("inviteID")
	if inviteID == "" {
		http.Error(w, "invite id is required", http.StatusBadRequest)
		return
	}

	if err := repositories.UpdateInviteStatus(h.DB, inviteID, userID, "declined"); err != nil {
		if err.Error() == "unauthorized" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if err.Error() == "invite not found" {
			http.Error(w, "invite not found", http.StatusNotFound)
			return
		}
		http.Error(w, "could not decline invite", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "invite declined"}`))
}

func (h *GroupHandler) RequestToJoin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getSessionUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	groupID := r.PathValue("groupID")
	if groupID == "" {
		http.Error(w, "group id is required", http.StatusBadRequest)
		return
	}

	req, err := repositories.CreateJoinRequest(h.DB, groupID, userID)
	if err != nil {
		if err.Error() == "already a member" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "could not create join request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func (h *GroupHandler) ApproveRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getSessionUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	requestID := r.PathValue("requestID")
	if requestID == "" {
		http.Error(w, "request id is required", http.StatusBadRequest)
		return
	}

	if err := repositories.UpdateRequestStatus(h.DB, requestID, userID, "approved"); err != nil {
		if err.Error() == "only the group creator can approve or decline requests" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if err.Error() == "request not found" {
			http.Error(w, "request not found", http.StatusNotFound)
			return
		}
		http.Error(w, "could not approve request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "request approved"}`))
}

func (h *GroupHandler) DeclineRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getSessionUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	requestID := r.PathValue("requestID")
	if requestID == "" {
		http.Error(w, "request id is required", http.StatusBadRequest)
		return
	}

	if err := repositories.UpdateRequestStatus(h.DB, requestID, userID, "declined"); err != nil {
		if err.Error() == "only the group creator can approve or decline requests" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if err.Error() == "request not found" {
			http.Error(w, "request not found", http.StatusNotFound)
			return
		}
		http.Error(w, "could not decline request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "request declined"}`))
}

func (h *GroupHandler) CreateGroupPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getSessionUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	groupID := r.PathValue("groupID")
	if groupID == "" {
		http.Error(w, "group id is required", http.StatusBadRequest)
		return
	}

	var body struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if body.Content == "" {
		http.Error(w, "content is required", http.StatusBadRequest)
		return
	}

	post, err := repositories.CreateGroupPost(h.DB, groupID, userID, body.Content)
	if err != nil {
		if err.Error() == "only group members can create posts" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "could not create post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// GetGroupFeed handles GET /api/groups/{groupID}/posts
func (h *GroupHandler) GetGroupFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getSessionUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	groupID := r.PathValue("groupID")
	if groupID == "" {
		http.Error(w, "group id is required", http.StatusBadRequest)
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 20

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	posts, err := repositories.GetGroupPosts(h.DB, groupID, userID, page, limit)
	if err != nil {
		if err.Error() == "only group members can view posts" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		http.Error(w, "could not fetch posts", http.StatusInternalServerError)
		return
	}

	if posts == nil {
		posts = []models.GroupPostWithMeta{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// CreateGroupComment handles POST /api/groups/{groupID}/posts/{postID}/comments
func (h *GroupHandler) CreateGroupComment(w http.ResponseWriter, r *http.Request) {
	// Gate 1 — method check
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Gate 2 — authentication
	userID, err := getSessionUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Gate 3 — extract URL values
	groupID := r.PathValue("groupID")
	postID := r.PathValue("postID")
	if groupID == "" || postID == "" {
		http.Error(w, "group id and post id are required", http.StatusBadRequest)
		return
	}

	// Gate 4 — decode body
	var body struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Gate 5 — validate content
	if body.Content == "" {
		http.Error(w, "content is required", http.StatusBadRequest)
		return
	}

	// Gate 6 — call repository
	comment, err := repositories.CreateGroupComment(h.DB, groupID, postID, userID, body.Content)
	if err != nil {
		if err.Error() == "only group members can comment" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if err.Error() == "post not found in this group" {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}
		http.Error(w, "could not create comment", http.StatusInternalServerError)
		return
	}

	// Gate 7 — send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

// GetGroupComments handles GET /api/groups/{groupID}/posts/{postID}/comments
func (h *GroupHandler) GetGroupComments(w http.ResponseWriter, r *http.Request) {
	// Gate 1 — method check
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Gate 2 — authentication
	userID, err := getSessionUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Gate 3 — extract URL values
	groupID := r.PathValue("groupID")
	postID := r.PathValue("postID")
	if groupID == "" || postID == "" {
		http.Error(w, "group id and post id are required", http.StatusBadRequest)
		return
	}

	// Gate 4 — call repository
	comments, err := repositories.GetGroupComments(h.DB, groupID, postID, userID)
	if err != nil {
		if err.Error() == "only group members can view comments" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if err.Error() == "post not found in this group" {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}
		http.Error(w, "could not fetch comments", http.StatusInternalServerError)
		return
	}

	// Gate 5 — send response
	if comments == nil {
		comments = []models.GroupCommentWithMeta{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func (h *GroupHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getSessionUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	groupID := r.PathValue("groupID")

	if groupID == "" {
		http.Error(w, "group id is required", http.StatusBadRequest)
		return
	}
	var req models.CreateEventRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	if req.Description==""{
		http.Error(w,"description is required",http.StatusBadRequest)
		return 
	}

	if req.EventDate==0{
		http.Error(w,"event date is required",http.StatusBadRequest)
		return
	}

	event,err:=repositories.CreateEvent(h.DB,groupID,userID,req)

	if err !=nil{
		if err.Error()=="only group members can create events"{
			http.Error(w,"forbidden",http.StatusForbidden)
		}
		http.Error(w,"could not create event",http.StatusInternalServerError)
	return
	}
w.Header().Set("Content-Type","application/json")
w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(event)
}

func (h *GroupHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getSessionUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	groupID := r.PathValue("groupID")
	if groupID == "" {
		http.Error(w, "group id is required", http.StatusBadRequest)
		return
	}

	events, err := repositories.GetGroupEvents(h.DB, groupID, userID)
	if err != nil {
		if err.Error() == "only group members can view events" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		http.Error(w, "could not fetch events", http.StatusInternalServerError)
		return
	}

	if events == nil {
		events = []models.EventWithResponses{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (h *GroupHandler) RespondToEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getSessionUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	eventID := r.PathValue("eventID")
	if eventID == "" {
		http.Error(w, "event id is required", http.StatusBadRequest)
		return
	}

	var body struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if body.Response != "going" && body.Response != "not_going" {
		http.Error(w, "response must be going or not_going", http.StatusBadRequest)
		return
	}

	resp, err := repositories.RespondToEvent(h.DB, eventID, userID, body.Response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *GroupHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getSessionUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	groupID := r.PathValue("groupID")
	if groupID == "" {
		http.Error(w, "missing group id", http.StatusBadRequest)
		return
	}

	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updated, err := repositories.UpdateGroup(h.DB, groupID, userID, body.Title, body.Description)
	if err != nil {
		if err.Error() == "only the group creator can edit the group" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if err.Error() == "title is required" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "could not update group", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func (h *GroupHandler) GetGroupMembersList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, err := getSessionUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	groupID := r.PathValue("groupID")
	if groupID == "" {
		http.Error(w, "missing group id", http.StatusBadRequest)
		return
	}

	type memberRow struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		AvatarURL string `json:"avatar_url"`
		JoinedAt  int64  `json:"joined_at"`
	}

	rows, err := h.DB.Query(`
		SELECT u.id, u.username, COALESCE(u.first_name,''), COALESCE(u.last_name,''), COALESCE(u.avatar_url,''), gm.joined_at
		FROM group_members gm
		JOIN users u ON u.id = gm.user_id
		WHERE gm.group_id = ?
		ORDER BY gm.joined_at ASC`, groupID)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var members []memberRow
	for rows.Next() {
		var m memberRow
		if err := rows.Scan(&m.ID, &m.Username, &m.FirstName, &m.LastName, &m.AvatarURL, &m.JoinedAt); err != nil {
			continue
		}
		members = append(members, m)
	}
	if members == nil {
		members = []memberRow{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}
