package repositories

import (
	"database/sql"
	"errors"
	"time"

	"social-network/backend/pkg/models"
)

type UpdateUserProfileParams struct {
	Username  *string
	FirstName *string
	LastName  *string
	Bio       *string
	AvatarURL *string
	Nickname  *string
}

func UpdateUserProfile(db *sql.DB, userID int, p UpdateUserProfileParams) (*models.User, error) {
	// Ensure at least one field is being updated.
	if p.Username == nil && p.FirstName == nil && p.LastName == nil && p.Bio == nil && p.AvatarURL == nil && p.Nickname == nil {
		return nil, errors.New("no fields to update")
	}

	fields := make([]string, 0, 6)
	args := make([]any, 0, 7)

	if p.Username != nil {
		fields = append(fields, "username = ?")
		args = append(args, *p.Username)
	}
	if p.FirstName != nil {
		fields = append(fields, "first_name = ?")
		args = append(args, *p.FirstName)
	}
	if p.LastName != nil {
		fields = append(fields, "last_name = ?")
		args = append(args, *p.LastName)
	}
	if p.Bio != nil {
		fields = append(fields, "bio = ?")
		args = append(args, *p.Bio)
	}
	if p.AvatarURL != nil {
		fields = append(fields, "avatar_url = ?")
		args = append(args, *p.AvatarURL)
	}
	if p.Nickname != nil {
		fields = append(fields, "nickname = ?")
		args = append(args, *p.Nickname)
	}

	// updated_at always updates.
	fields = append(fields, "updated_at = CURRENT_TIMESTAMP")

	args = append(args, userID)

	q := "UPDATE users SET " + join(fields, ", ") + " WHERE id = ?"
	res, err := db.Exec(q, args...)
	if err != nil {
		return nil, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, ErrUserNotFound
	}

	return GetUserByID(db, userID)
}

func SetUserPrivacy(db *sql.DB, userID int, isPrivate bool) (*models.User, error) {
	res, err := db.Exec(`UPDATE users SET is_private = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, isPrivate, userID)
	if err != nil {
		return nil, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, ErrUserNotFound
	}
	return GetUserByID(db, userID)
}

// small helper to avoid importing strings everywhere
func join(parts []string, sep string) string {
	if len(parts) == 0 {
		return ""
	}
	out := parts[0]
	for i := 1; i < len(parts); i++ {
		out += sep + parts[i]
	}
	// silence unused import warning if time isn't used elsewhere in file
	_ = time.Now
	return out
}
