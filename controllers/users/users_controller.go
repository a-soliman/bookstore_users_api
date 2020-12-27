package users

import (
	"net/http"
	"strconv"

	"github.com/a-soliman/bookstore_users_api/domain/users"
	"github.com/a-soliman/bookstore_users_api/services"
	"github.com/a-soliman/bookstore_users_api/utils/errors"
	"github.com/gin-gonic/gin"
)

// GetUser finds a user by id
func GetUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		userErr := errors.NewBadRequestError("invalid user id")
		c.JSON(userErr.Status, userErr)
		return
	}
	result, userErr := services.GetUser(userID)
	if userErr != nil {
		c.JSON(userErr.Status, userErr)
		return
	}
	c.JSON(http.StatusOK, result)
}

// CreateUser creates a new user
func CreateUser(c *gin.Context) {
	var user users.User

	// read the body from request, and bind user JSON
	if err := c.ShouldBindJSON(&user); err != nil {
		jsonErr := errors.NewBadRequestError("invalid json body")
		c.JSON(jsonErr.Status, jsonErr)
		return
	}

	result, saveErr := services.CreateUser(&user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}

// UpdateUser updates a user entity
func UpdateUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		userErr := errors.NewBadRequestError("invalid user id")
		c.JSON(userErr.Status, userErr)
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

	result, updateErr := services.UpdateUser(isPartial, &user)
	if updateErr != nil {
		c.JSON(updateErr.Status, updateErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}

// DeleteUser delete a given user
func DeleteUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		userErr := errors.NewBadRequestError("invalid user id")
		c.JSON(userErr.Status, userErr)
		return
	}

	user := users.User{ID: userID}

	result, deleteErr := services.DeleteUser(&user)
	if deleteErr != nil {
		c.JSON(deleteErr.Status, deleteErr)
		return
	}
	c.JSON(http.StatusOK, result)
}
