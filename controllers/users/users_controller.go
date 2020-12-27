package users

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-soliman/bookstore_users_api/domain/users"
	"github.com/a-soliman/bookstore_users_api/services"
	"github.com/a-soliman/bookstore_users_api/utils/errors"
	"github.com/gin-gonic/gin"
)

func getUserID(userIDParam string) (int64, *errors.RestErr) {
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		return 0, errors.NewBadRequestError("invalid user id")
	}
	return userID, nil
}

// Get finds a user by id
func Get(c *gin.Context) {
	userID, idErr := getUserID(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}
	result, userErr := services.UsersService.GetUser(userID)
	if userErr != nil {
		c.JSON(userErr.Status, userErr)
		return
	}
	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-PUBLIC") == "true"))
}

// Search finds a list of users by given search query
func Search(c *gin.Context) {
	status := strings.TrimSpace(c.Query("status"))
	fmt.Println(status)
	if status == "" {
		err := errors.NewBadRequestError("missing status")
		c.JSON(err.Status, err)
		return
	}
	result, findErr := services.UsersService.SearchUser(status)
	if findErr != nil {
		c.JSON(findErr.Status, findErr)
		return
	}
	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-PUBLIC") == "true"))
}

// Create creates a new user
func Create(c *gin.Context) {
	var user users.User

	// read the body from request, and bind user JSON
	if err := c.ShouldBindJSON(&user); err != nil {
		jsonErr := errors.NewBadRequestError("invalid json body")
		c.JSON(jsonErr.Status, jsonErr)
		return
	}

	result, saveErr := services.UsersService.CreateUser(&user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-PUBLIC") == "true"))
}

// Update updates a user entity
func Update(c *gin.Context) {
	userID, idErr := getUserID(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	var user users.User

	// read the body from request, and bind user JSON
	if err := c.ShouldBindJSON(&user); err != nil {
		jsonErr := errors.NewBadRequestError("invalid json body")
		c.JSON(jsonErr.Status, jsonErr)
		return
	}

	user.ID = userID

	isPartial := c.Request.Method == http.MethodPatch

	result, updateErr := services.UsersService.UpdateUser(isPartial, &user)
	if updateErr != nil {
		c.JSON(updateErr.Status, updateErr)
		return
	}
	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-PUBLIC") == "true"))
}

// Delete delete a given user
func Delete(c *gin.Context) {
	userID, idErr := getUserID(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	user := users.User{ID: userID}

	result, deleteErr := services.UsersService.DeleteUser(&user)
	if deleteErr != nil {
		c.JSON(deleteErr.Status, deleteErr)
		return
	}
	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-PUBLIC") == "true"))
}
