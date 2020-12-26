package users

import (
	"fmt"

	"github.com/a-soliman/bookstore_users_api/utils/dates"
	"github.com/a-soliman/bookstore_users_api/utils/errors"
)

var (
	usersDB = make(map[int64]*User)
)

// Get fetches a user by userId from database
func (u *User) Get() *errors.RestErr {
	result, exists := usersDB[u.ID]
	if !exists {
		return errors.NewNotFoundError(fmt.Sprintf("user %d not found", u.ID))
	}
	u.FirstName = result.FirstName
	u.LastName = result.LastName
	u.Email = result.Email
	u.CreatedAt = result.CreatedAt

	return nil
}

// Save persists a given user into the database
func (u *User) Save() *errors.RestErr {
	// check if already exists
	existingUser, exists := usersDB[u.ID]
	if exists {
		if existingUser.Email == u.Email {
			return errors.NewBadRequestError(fmt.Sprintf("email %s already registered", u.Email))
		}
		return errors.NewBadRequestError(fmt.Sprintf("user %d already exists", u.ID))
	}
	// append UTC time for createdAt
	u.CreatedAt = dates.GetNowString()
	// save the user
	usersDB[u.ID] = u
	return nil
}
