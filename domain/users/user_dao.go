package users

import (
	"fmt"
	"strings"

	"github.com/a-soliman/bookstore_users_api/datasources/mysql/users_db"
	"github.com/a-soliman/bookstore_users_api/utils/dates"
	"github.com/a-soliman/bookstore_users_api/utils/errors"
)

const (
	indexUniqueEmail = "email_UNIQUE"
	errorNoRows      = "no rows in result set"
	queryInsertUser  = "INSERT INTO users(first_name, last_name, email, created_at) VALUES(?, ?, ?, ?);"
	queryGetUser     = "SELECT id, first_name, last_name, email, created_at FROM users WHERE id=?;"
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
	if err := result.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.CreatedAt); err != nil {
		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundError(fmt.Sprintf("user %d not found", u.ID))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error while trying to get user by id, id = %d: err = %s", u.ID, err.Error()))
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

	insertResult, err := stmt.Exec(u.FirstName, u.LastName, u.Email, u.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), indexUniqueEmail) {
			return errors.NewBadRequestError(fmt.Sprintf("email %s already exists", u.Email))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error while trying to save user: %s", err.Error()))
	}

	userID, err := insertResult.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error while trying to save user: %s", err.Error()))
	}
	u.ID = userID
	return nil
}
