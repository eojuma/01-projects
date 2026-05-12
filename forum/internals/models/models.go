package models

type ResponseData struct {
	IsLoggedIn bool
	Posts      []Post
	Likes      []Like
	Comments   []Comment
	User       User
	Initials   string
	StatusCode int
	Error string

}

type User struct {
	UserID         string
	Email          string
	Username       string
	Password       string
	ProfilePicture string
}

type Post struct {
	PostID    string
	UserID    string
	Title     string
	Content   string
	Category  string
	CreatedAt string
	Image     string
	Username  string
	Initials  string

	Likes int
	Dislikes int
}

type Comment struct {
	CommentID string
	UserID    string
	PostID    string
	Content   string
	CreatedAt string
	Username  string
	Likes     int
	Dislikes  int
}

type Like struct {
	LikeID string
	Counts int
	UserID string
	PostID string
}

type Dislike struct {
	DislikeID string
	Counts    int
	UserID   string
	PostID   string
}

type CommentLike struct {
	LikeID    string
	Counts    int
	UserID    string
	CommentID string
}

type CommentDislike struct {
	DislikeID    string
	Counts       int
	UserID       string
	CommentID    string
}

type Session struct {
	SessionID string
	UserID    string
	ExpiresAt string
}

var Data = ResponseData{
	IsLoggedIn: false,
}

var CurrentUser = User{}
