package main

import (
	"fmt"
	"net/http"

	DB "forum/internals/db"
	"forum/internals/handlers"
	"forum/internals/middleware"
)

func main() {
	fmt.Println("Starting server...")

	fmt.Println("Initializing DB...")
	db, err := DB.InitDB("./internals/db/forum.db")
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	defer db.Close()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))
	http.HandleFunc("/login", handlers.LoginHandler(db))
	http.HandleFunc("/register", handlers.RegisterHandler(db))
	http.HandleFunc("/logout", handlers.Logout(db))
	http.HandleFunc("/search", handlers.SearchHandler(db))
	http.HandleFunc("/post/create", handlers.CreatePostHandler(db))
	http.HandleFunc("/post/category", handlers.FilterByCategory(db))
	http.HandleFunc("/post/my-likes", handlers.MyLikesHandler(db))
	http.HandleFunc("/post/my-posts", handlers.MyPostsHandler(db))
	http.HandleFunc("/post/likes", handlers.SortByLikes(db))
	http.HandleFunc("/post/comments", handlers.SortByComments(db))
	http.HandleFunc("/post/view", handlers.ViewPostHandler(db))
	http.HandleFunc("/comment/vote/like", middleware.Authorize(db, handlers.CommentLikeHandler(db)))
	http.HandleFunc("/comment/vote/dislike", middleware.Authorize(db, handlers.CommentDislikeHandler(db)))
	http.HandleFunc("/post-comment", middleware.Authorize(db, handlers.CommentHandler(db)))
	http.HandleFunc("/vote/like", middleware.Authorize(db, handlers.LikeHandler(db)))
	http.HandleFunc("/vote/dislike", middleware.Authorize(db, handlers.DislikeHandler(db)))
	http.HandleFunc("/", handlers.HomeHandler(db))

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server Error: ", err)
	}
}
