package users

import (
	"fmt"

	"github.com/a-soliman/bookstore_users_api/datasources/mysql/users_db"
	"github.com/a-soliman/bookstore_users_api/logger"
	"github.com/a-soliman/bookstore_users_api/utils/errors"
)

const (
	queryInsertUser       = "INSERT INTO users(first_name, last_name, email, created_at, password, status) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser          = "SELECT id, first_name, last_name, email, created_at, status FROM users WHERE id=?;"
	queryUpdateUser       = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser       = "DELETE FROM users WHERE id=?;"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, created_at, status FROM users WHERE status=?"
)

var (
	usersDB = make(map[int64]*User)
)

// Get fetches a user by userId from database
func (u *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error while trying to prepare get user statement", err)
		return errors.NewInternalServerError(errors.DatabaseErrorMsg)
	}
	defer stmt.Close()

	result := stmt.QueryRow(u.ID)
	if getErr := result.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.CreatedAt, &u.Status); getErr != nil {
		logger.Error("error while trying to get user by id", getErr)
		return errors.NewInternalServerError(errors.DatabaseErrorMsg)
	}

	return nil
}

// FindByStatus given a status string it returns a slice of users and a restErr
func (u *User) FindByStatus(status string) (*Users, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("error while trying to prepare find user statement", err)
		return nil, errors.NewInternalServerError(errors.DatabaseErrorMsg)
	}
	defer stmt.Close()

	rows, findErr := stmt.Query(status)
	if findErr != nil {
		logger.Error("error while trying to find user by status", findErr)
		return nil, errors.NewInternalServerError(errors.DatabaseErrorMsg)
	}
	defer rows.Close()

	results := Users{}
	for rows.Next() {
		var user *User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.Status); err != nil {
			logger.Error("error while trying to scan user in find user by status", err)
			return nil, errors.NewInternalServerError(errors.DatabaseErrorMsg)
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return &results, nil
}

// Save persists a given user into the database
func (u *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error while trying to prepare save user statement", err)
		return errors.NewInternalServerError(errors.DatabaseErrorMsg)
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(u.FirstName, u.LastName, u.Email, u.CreatedAt, u.Password, u.Status)
	if saveErr != nil {
		logger.Error("error while trying to save user", saveErr)
		return errors.NewInternalServerError(errors.DatabaseErrorMsg)
	}

	userID, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error while trying to save user, LastInsertId", err)
		return errors.NewInternalServerError(errors.DatabaseErrorMsg)
	}
	u.ID = userID
	return nil
}

// Update updates the user entity
func (u *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error while trying to prepare update user statement", err)
		return errors.NewInternalServerError(errors.DatabaseErrorMsg)
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(u.FirstName, u.LastName, u.Email, u.ID)
	if updateErr != nil {
		logger.Error("error while trying to update user", updateErr)
		return errors.NewInternalServerError(errors.DatabaseErrorMsg)
	}
	return nil
}

// Delete deletes a user from DB
func (u *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error while trying to prepare delete user statement", err)
		return errors.NewInternalServerError(errors.DatabaseErrorMsg)
	}
	defer stmt.Close()
	_, deleteErr := stmt.Exec(u.ID)
	if deleteErr != nil {
		logger.Error("error while trying to delete user", deleteErr)
		return errors.NewInternalServerError(errors.DatabaseErrorMsg)
	}
	return nil
}
