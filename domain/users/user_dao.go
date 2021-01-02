package users

import (
	"fmt"
	"strings"

	"github.com/a-soliman/bookstore_oauth-go/oauth/errors"
	"github.com/a-soliman/bookstore_users_api/datasources/mysql/users_db"
	"github.com/a-soliman/bookstore_users_api/logger"
	"github.com/a-soliman/bookstore_utils-go/rest_errors"
)

const (
	queryInsertUser            = "INSERT INTO users(first_name, last_name, email, created_at, password, status) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser               = "SELECT id, first_name, last_name, email, created_at, status FROM users WHERE id=?;"
	queryUpdateUser            = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser            = "DELETE FROM users WHERE id=?;"
	queryFindByStatus          = "SELECT id, first_name, last_name, email, created_at, status FROM users WHERE status=?"
	queryFindByEmailAndPasword = "SELECT id, first_name, last_name, email, created_at, status FROM users WHERE email=? AND password=? AND status=?;"
)

var (
	usersDB = make(map[int64]*User)
)

// Get fetches a user by userId from database
func (u *User) Get() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error while trying to prepare get user statement", err)
		return rest_errors.NewInternalServerError(errors.DatabaseErrorMsg, err)
	}
	defer stmt.Close()

	result := stmt.QueryRow(u.ID)
	if getErr := result.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.CreatedAt, &u.Status); getErr != nil {
		logger.Error("error while trying to get user by id", getErr)
		return rest_errors.NewInternalServerError(errors.DatabaseErrorMsg, err)
	}

	return nil
}

// FindByStatus given a status string it returns a slice of users and a restErr
func (u *User) FindByStatus(status string) (*Users, *rest_errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindByStatus)
	if err != nil {
		logger.Error("error while trying to prepare find user statement", err)
		return nil, rest_errors.NewInternalServerError(errors.DatabaseErrorMsg, err)
	}
	defer stmt.Close()

	rows, findErr := stmt.Query(status)
	if findErr != nil {
		logger.Error("error while trying to find user by status", findErr)
		return nil, rest_errors.NewInternalServerError(errors.DatabaseErrorMsg, err)
	}
	defer rows.Close()

	results := Users{}
	for rows.Next() {
		var user *User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.Status); err != nil {
			logger.Error("error while trying to scan user in find user by status", err)
			return nil, rest_errors.NewInternalServerError(errors.DatabaseErrorMsg, err)
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return &results, nil
}

// FindByEmailAndPassword fetches a user by email and password from database
func (u *User) FindByEmailAndPassword() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryFindByEmailAndPasword)
	if err != nil {
		logger.Error("error while trying to prepare find user by email and password statement", err)
		return rest_errors.NewInternalServerError(errors.DatabaseErrorMsg, err)
	}
	defer stmt.Close()

	result := stmt.QueryRow(u.Email, u.Password, StatusActive)
	if getErr := result.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.CreatedAt, &u.Status); getErr != nil {
		if strings.Contains(getErr.Error(), "no rows in result set") {
			return rest_errors.NewNotFoundError("invalid user credentials")
		}
		logger.Error("error while trying to find user by email and password", getErr)
		return rest_errors.NewInternalServerError(errors.DatabaseErrorMsg, getErr)
	}

	return nil
}

// Save persists a given user into the database
func (u *User) Save() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error while trying to prepare save user statement", err)
		return rest_errors.NewInternalServerError(errors.DatabaseErrorMsg, err)
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(u.FirstName, u.LastName, u.Email, u.CreatedAt, u.Password, u.Status)
	if saveErr != nil {
		logger.Error("error while trying to save user", saveErr)
		return rest_errors.NewInternalServerError(errors.DatabaseErrorMsg, saveErr)
	}

	userID, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error while trying to save user, LastInsertId", err)
		return rest_errors.NewInternalServerError(errors.DatabaseErrorMsg, err)
	}
	u.ID = userID
	return nil
}

// Update updates the user entity
func (u *User) Update() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error while trying to prepare update user statement", err)
		return rest_errors.NewInternalServerError(errors.DatabaseErrorMsg, err)
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(u.FirstName, u.LastName, u.Email, u.ID)
	if updateErr != nil {
		logger.Error("error while trying to update user", updateErr)
		return rest_errors.NewInternalServerError(errors.DatabaseErrorMsg, updateErr)
	}
	return nil
}

// Delete deletes a user from DB
func (u *User) Delete() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error while trying to prepare delete user statement", err)
		return rest_errors.NewInternalServerError(errors.DatabaseErrorMsg, err)
	}
	defer stmt.Close()
	_, deleteErr := stmt.Exec(u.ID)
	if deleteErr != nil {
		logger.Error("error while trying to delete user", deleteErr)
		return rest_errors.NewInternalServerError(errors.DatabaseErrorMsg, deleteErr)
	}
	return nil
}
