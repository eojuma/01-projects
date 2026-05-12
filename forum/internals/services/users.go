package services

import (
	"database/sql"
	"fmt"

	"forum/internals/models"
)

func AddUser(db *sql.DB, user models.User) error {
	addUser, err := db.Prepare("INSERT INTO USERS VALUES(?,?,?,?,?)")
	if err != nil {
		fmt.Println("Error adding user: ", err)
		return err
	}
	defer addUser.Close()
	_, err = addUser.Exec(user.UserID, user.Email, user.Username, user.Password, user.ProfilePicture)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(db *sql.DB, email, username string) (*models.User, error) {
	user := &models.User{}
	var viewUsers *sql.Stmt
	viewUsers, err := db.Prepare("SELECT * FROM USERS WHERE email=? OR username=?")
	if err != nil {
		return nil, err
	}
	defer viewUsers.Close()
	users, err := viewUsers.Query(email, username)
	if err != nil {
		return nil, err
	}
	for users.Next() {
		err = users.Scan(&user.UserID, &user.Email, &user.Username, &user.Password, &user.ProfilePicture)
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}

func EditUsers(db *sql.DB, user models.User) {
	editUser, err := db.Prepare("UPDATE USERS SET userId=?, email=?, username=?,password=?,profile=? WHERE userId=?")
	if err != nil {
		fmt.Println("Error updating user: ", err)
		return
	}
	defer editUser.Close()
	users, err := editUser.Exec()
	if err != nil {
		fmt.Println("Error retrieving users: ", err)
		return
	}
	fmt.Println(users.RowsAffected())
}

// End User
