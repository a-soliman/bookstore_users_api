package users

import (
	"strings"

	"github.com/a-soliman/bookstore_utils-go/rest_errors"
)

const (
	// StatusActive the default status users
	StatusActive = "active"
)

// User the user entity
type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
	Password  string `json:"password"`
}

// Users a slice of user pointer
type Users []*User

// Validate returns a pointer to an error if the given user is invalid
func (u *User) Validate() rest_errors.RestErr {
	u.FirstName = strings.TrimSpace(u.FirstName)
	if u.FirstName == "" {
		return rest_errors.NewBadRequestError("firstname is required")
	}
	u.LastName = strings.TrimSpace(u.LastName)
	if u.LastName == "" {
		return rest_errors.NewBadRequestError("lastname is required")
	}
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	if u.Email == "" {
		return rest_errors.NewBadRequestError("invalid email address")
	}
	u.Password = strings.TrimSpace(u.Password)
	if u.Password == "" {
		return rest_errors.NewBadRequestError("invalid password")
	}
	return nil
}
