package app

import (
	"github.com/a-soliman/bookstore_users_api/logger"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

// StartApplication starts the api
func StartApplication() {
	mapUrls()
	logger.Info("about to start the application...")
	router.Run(":8082")
}
