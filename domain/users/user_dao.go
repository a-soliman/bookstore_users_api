package users

import (
	"fmt"

	"github.com/a-soliman/bookstore_users_api/datasources/mysql/users_db"
	"github.com/a-soliman/bookstore_users_api/utils/dates"
	"github.com/a-soliman/bookstore_users_api/utils/errors"
	"github.com/a-soliman/bookstore_users_api/utils/mysql_utils"
)

const (
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, created_at) VALUES(?, ?, ?, ?);"
	queryGetUser    = "SELECT id, first_name, last_name, email, created_at FROM users WHERE id=?;"
	queryUpdateUser = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser = "DELETE FROM users WHERE id=?;"
)

var (
	usersDB = make(map[int64]*User)
)

// Get fetches a user by userId from database
func (u *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	result := stmt.QueryRow(u.ID)
	if getErr := result.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.CreatedAt); getErr != nil {
		return mysql_utils.ParseError(getErr)
	}

	return nil
}

// Save persists a given user into the database
func (u *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	// append UTC time for createdAt
	u.CreatedAt = dates.GetNowString()

	insertResult, saveErr := stmt.Exec(u.FirstName, u.LastName, u.Email, u.CreatedAt)
	if saveErr != nil {
		return mysql_utils.ParseError(saveErr)
	}

	userID, err := insertResult.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error while trying to save user: %s", err.Error()))
	}
	u.ID = userID
	return nil
}

// Update updates the user entity
func (u *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(u.FirstName, u.LastName, u.Email, u.ID)
	if updateErr != nil {
		return mysql_utils.ParseError(updateErr)
	}
	return nil
}

// Delete deletes a user from DB
func (u *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	_, deleteErr := stmt.Exec(u.ID)
	if deleteErr != nil {
		return mysql_utils.ParseError(deleteErr)
	}
	return nil
}
