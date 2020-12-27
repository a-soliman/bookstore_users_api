package services

import (
	"github.com/a-soliman/bookstore_users_api/domain/users"
	"github.com/a-soliman/bookstore_users_api/utils/crypto_utils"
	"github.com/a-soliman/bookstore_users_api/utils/dates"
	"github.com/a-soliman/bookstore_users_api/utils/errors"
)

// GetUser gets user by id
func GetUser(userID int64) (*users.User, *errors.RestErr) {
	user := &users.User{ID: userID}
	if err := user.Get(); err != nil {
		return nil, err
	}
	return user, nil
}

// Search returns a list of users and an error
func Search(status string) (*users.Users, *errors.RestErr) {
	var dao users.User
	return dao.FindByStatus(status)
}

// CreateUser creates a new user
func CreateUser(user *users.User) (*users.User, *errors.RestErr) {
	// validate user
	if err := user.Validate(); err != nil {
		return nil, err
	}

	user.Status = users.StatusActive
	user.CreatedAt = dates.GetNowDBFormat()
	user.Password = crypto_utils.GetMd5(user.Password)

	// attempt to save the user
	if err := user.Save(); err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser updates a user entity
func UpdateUser(isPartial bool, user *users.User) (*users.User, *errors.RestErr) {
	// fetch current user entity
	current, err := GetUser(user.ID)
	if err != nil {
		return nil, err
	}

	if isPartial {
		if user.FirstName != "" {
			current.FirstName = user.FirstName
		}
		if user.LastName != "" {
			current.LastName = user.LastName
		}
		if err := user.Validate(); err == nil {
			current.Email = user.Email
		}
	} else {
		current.FirstName = user.FirstName
		current.LastName = user.LastName
		current.Email = user.Email
	}

	// attempt to update the user
	if err := user.Update(); err != nil {
		return nil, err
	}
	return current, nil
}

// DeleteUser deletes a user from DB and returns the deletedUser or err
func DeleteUser(user *users.User) (*users.User, *errors.RestErr) {
	user, err := GetUser(user.ID)
	if err != nil {
		return nil, err
	}

	if err := user.Delete(); err != nil {
		return nil, err
	}
	return user, nil
}
