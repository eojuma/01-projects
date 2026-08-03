package models

type Group struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatorID   string `json:"creator_id"`
	CreatedAt   int64  `json:"created_at"`
}

type CreateGroupRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GroupMember struct {
	GroupID  string `json:"group_id"`
	UserID   string `json:"user_id"`
	JoinedAt int64  `json:"joined_at"`
}

type GroupPost struct {
	ID        string `json:"id"`
	GroupID   string `json:"group_id"`
	AuthorID  string `json:"author_id"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
}

type GroupComment struct {
	ID        string `json:"id"`
	PostID    string `json:"post_id"`
	AuthorID  string `json:"author_id"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
}

type Event struct {
	ID          string `json:"id"`
	GroupID     string `json:"group_id"`
	CreatorID   string `json:"creator_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	EventDate   int64  `json:"event_date"`
	CreatedAt   int64  `json:"created_at"`
}

type CreateEventRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	EventDate   int64  `json:"event_date"`
}

type EventResponse struct {
	ID        string `json:"id"`
	EventID   string `json:"event_id"`
	UserID    string `json:"user_id"`
	Response  string `json:"response"`
	CreatedAt int64  `json:"created_at"`
}

type GroupWithCount struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatorID   string `json:"creator_id"`
	CreatedAt   int64  `json:"created_at"`
	MemberCount int    `json:"member_count"`
	IsMember    bool   `json:"is_member"`
	IsInvited   bool   `json:"is_invited"`
	InviteID    string `json:"invite_id,omitempty"`
	IsRequested bool   `json:"is_requested"`
}

type GroupInvite struct {
	ID        string `json:"id"`
	GroupID   string `json:"group_id"`
	InviterID string `json:"inviter_id"`
	InviteeID string `json:"invitee_id"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
}

type GroupJoinRequest struct {
	ID        string `json:"id"`
	GroupID   string `json:"group_id"`
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
}

type GroupPostWithMeta struct {
	ID             string `json:"id"`
	GroupID        string `json:"group_id"`
	AuthorID       string `json:"author_id"`
	AuthorUsername string `json:"author_username"`
	Content        string `json:"content"`
	CreatedAt      int64  `json:"created_at"`
	CommentCount   int    `json:"comment_count"`
}

type GroupCommentWithMeta struct {
	ID             string `json:"id"`
	PostID         string `json:"post_id"`
	AuthorID       string `json:"author_id"`
	AuthorUsername string `json:"author_username"`
	Content        string `json:"content"`
	CreatedAt      int64  `json:"created_at"`
}

type EventWithResponses struct {
	ID          string          `json:"id"`
	GroupID     string          `json:"group_id"`
	CreatorID   string          `json:"creator_id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	EventDate   int64           `json:"event_date"`
	CreatedAt   int64           `json:"created_at"`
	Going       int             `json:"going"`
	NotGoing    int             `json:"not_going"`
	UserResponse string         `json:"user_response"`
}

type GroupJoinRequestWithUser struct {
	GroupJoinRequest
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}