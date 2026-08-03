package models

type FollowResult struct {
	Following bool   `json:"following"`
	Status    string `json:"status"`
}

type NotificationData struct {
	NotificationID string `json:"notification_id"`
	SenderID       string `json:"sender_id"`
	Type           string `json:"type"`
}

type UserSummary struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	AvatarURL   string `json:"avatar_url"`
	IsFollowing bool   `json:"is_following"`
}

type PaginatedUsers struct {
	Users   []UserSummary `json:"users"`
	Total   int           `json:"total"`
	Page    int           `json:"page"`
	PerPage int           `json:"per_page"`
}
