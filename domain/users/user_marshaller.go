package users

import "encoding/json"

// PublicUser a struct for the public user, no password or email
type PublicUser struct {
	ID        int64  `json:"id"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
}

// PrivateUser a struct for internal user, where all fields are accessable (except the password)
type PrivateUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
}

// Marshall returns the appropriate user struct depends on where the request is public
func (user *User) Marshall(isPublic bool) interface{} {
	userJSON, _ := json.Marshal(user)
	var res interface{}
	if isPublic {
		var publicUser PublicUser
		json.Unmarshal(userJSON, &publicUser)
		res = publicUser
	} else {
		var privateUser PrivateUser
		json.Unmarshal(userJSON, &privateUser)
		res = privateUser
	}
	return res
}

// Marshall returns the appropriate user struct depends on where the request is public
func (users *Users) Marshall(isPublic bool) []interface{} {
	result := make([]interface{}, len(*users))
	for idx, user := range *users {
		result[idx] = user.Marshall(isPublic)
	}
	return result
}
