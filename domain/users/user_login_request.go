package users

// LoginRequest a struct
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
