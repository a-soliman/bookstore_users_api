package users

import (
	"fmt"

	"github.com/a-soliman/bookstore_users_api/datasources/mysql/users_db"
	"github.com/a-soliman/bookstore_users_api/utils/errors"
	"github.com/a-soliman/bookstore_users_api/utils/mysql_utils"
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
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	result := stmt.QueryRow(u.ID)
	if getErr := result.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.CreatedAt, &u.Status); getErr != nil {
		return mysql_utils.ParseError(getErr)
	}

	return nil
}

// FindByStatus given a status string it returns a slice of users and a restErr
func (u *User) FindByStatus(status string) (*Users, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	rows, findErr := stmt.Query(status)
	if findErr != nil {
		return nil, mysql_utils.ParseError(findErr)
	}
	defer rows.Close()

	results := Users{}
	for rows.Next() {
		var user *User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.Status); err != nil {
			return nil, mysql_utils.ParseError(findErr)
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
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(u.FirstName, u.LastName, u.Email, u.CreatedAt, u.Password, u.Status)
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
