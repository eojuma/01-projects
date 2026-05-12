package services

import (
	"database/sql"
	"forum/internals/models"
)

func AddDislike(db *sql.DB, dislike models.Dislike) error {
	dislikee, err := db.Prepare("INSERT INTO dislikes (dislikeId, counts, userId, postId) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer dislikee.Close()

	_, err = dislikee.Exec(dislike.DislikeID, dislike.Counts, dislike.UserID, dislike.PostID)
	if err != nil {
		return err
	}

	return nil
}

func RemoveDislike(db *sql.DB, userID string, postID string) error {
	_, err := db.Exec("DELETE FROM dislikes WHERE userID = ? AND postID = ?", userID, postID)
	if err != nil {
		return err
	}

	return nil
}

func HasUserDisliked(db *sql.DB, userID string, postID string) (bool, error) {
	var dislikedID string
	err := db.QueryRow("SELECT dislikeID FROM dislikes WHERE userID = ? and postID = ?", userID, postID).Scan(&dislikedID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}


func GetDislikeCount(db *sql.DB, postID string) (int, error) {
	var count int
	err := db.QueryRow("SELECT SUM(*) FROM dislikes WHERE postId = ?", postID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

