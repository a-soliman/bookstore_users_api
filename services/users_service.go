package services

import (
	"github.com/a-soliman/bookstore_users_api/domain/users"
	"github.com/a-soliman/bookstore_utils-go/crypto_utils"
	"github.com/a-soliman/bookstore_utils-go/date_utils"
	"github.com/a-soliman/bookstore_utils-go/rest_errors"
)

var (
	// UsersService the usersService
	UsersService userServiceInterface = &usersService{}
)

type usersService struct{}

type userServiceInterface interface {
	GetUser(int64) (*users.User, rest_errors.RestErr)
	SearchUser(string) (*users.Users, rest_errors.RestErr)
	CreateUser(*users.User) (*users.User, rest_errors.RestErr)
	UpdateUser(bool, *users.User) (*users.User, rest_errors.RestErr)
	DeleteUser(*users.User) (*users.User, rest_errors.RestErr)
	LoginUser(users.LoginRequest) (*users.User, rest_errors.RestErr)
}

// GetUser gets user by id
func (us *usersService) GetUser(userID int64) (*users.User, rest_errors.RestErr) {
	user := &users.User{ID: userID}
	if err := user.Get(); err != nil {
		return nil, err
	}
	return user, nil
}

// Search returns a list of users and an error
func (us *usersService) SearchUser(status string) (*users.Users, rest_errors.RestErr) {
	var dao users.User
	return dao.FindByStatus(status)
}

func (us *usersService) LoginUser(request users.LoginRequest) (*users.User, rest_errors.RestErr) {
	dao := &users.User{
		Email:    request.Email,
		Password: crypto_utils.GetMd5(request.Password),
	}

	if err := dao.FindByEmailAndPassword(); err != nil {
		return nil, err
	}
	return dao, nil
}

// CreateUser creates a new user
func (us *usersService) CreateUser(user *users.User) (*users.User, rest_errors.RestErr) {
	// validate user
	if err := user.Validate(); err != nil {
		return nil, err
	}

	user.Status = users.StatusActive
	user.CreatedAt = date_utils.GetNowDBFormat()
	user.Password = crypto_utils.GetMd5(user.Password)

	// attempt to save the user
	if err := user.Save(); err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser updates a user entity
func (us *usersService) UpdateUser(isPartial bool, user *users.User) (*users.User, rest_errors.RestErr) {
	// fetch current user entity
	current, err := us.GetUser(user.ID)
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
func (us *usersService) DeleteUser(user *users.User) (*users.User, rest_errors.RestErr) {
	user, err := us.GetUser(user.ID)
	if err != nil {
		return nil, err
	}

	if err := user.Delete(); err != nil {
		return nil, err
	}
	return user, nil
}
