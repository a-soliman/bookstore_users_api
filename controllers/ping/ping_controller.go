package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping handles Ping request
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
