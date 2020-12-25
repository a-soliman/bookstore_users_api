package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUser finds a user by id
func GetUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "implement me!")
}

// CreateUser creates a new user
func CreateUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "implement me!")
}
