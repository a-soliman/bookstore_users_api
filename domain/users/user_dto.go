package users

import (
	"strings"

	"github.com/a-soliman/bookstore_users_api/utils/errors"
)

// User the user entity
type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

// Validate returns a pointer to an error if the given user is invalid
func (u *User) Validate() *errors.RestErr {
	// email
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	if u.Email == "" {
		return errors.NewBadRequestError("invalid email address")
	}
	return nil
}
