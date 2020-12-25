package app

import "github.com/a-soliman/bookstore_users_api/controllers"

func mapUrls() {
	router.GET("/ping", controllers.Ping)
}
