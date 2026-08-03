package models

type Post struct {
	ID              string `json:"id"`
	AuthorID        string `json:"author_id"`
	AuthorUsername  string `json:"author_username,omitempty"`
	AuthorIsPrivate bool   `json:"author_is_private,omitempty"`
	Content         string `json:"content"`
	Privacy         string `json:"privacy"`
	CreatedAt       int64  `json:"created_at"`
}

type CreatePostInput struct {
	ID           string
	AuthorID     string
	Content      string
	Privacy      string
	ImagePaths   []string
	AllowedUsers []string
}

type FeedPost struct {
	ID              string   `json:"id"`
	AuthorID        string   `json:"author_id"`
	AuthorUsername  string   `json:"author_username,omitempty"`
	AuthorIsPrivate bool     `json:"author_is_private"`
	Content         string   `json:"content"`
	Privacy         string   `json:"privacy"`
	CreatedAt       int64    `json:"created_at"`
	ImageURLs       []string `json:"image_urls"`
	CommentCount    int      `json:"comment_count"`
	LikesCount      int      `json:"likes_count"`
	Liked           bool     `json:"liked"`
	FollowStatus    string   `json:"follow_status,omitempty"`
}

type FeedResponse struct {
	Posts      []FeedPost `json:"posts"`
	NextCursor string     `json:"next_cursor,omitempty"`
}

type Comment struct {
	ID        string      `json:"id"`
	PostID    string      `json:"post_id"`
	AuthorID  string      `json:"author_id"`
	Content   string      `json:"content"`
	ImagePath string      `json:"image_path,omitempty"`
	CreatedAt int64       `json:"created_at"`
	Author    UserSummary `json:"author"`
}

type CreateCommentInput struct {
	ID        string
	PostID    string
	AuthorID  string
	Content   string
	ImagePath string
}
