package users

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-soliman/bookstore_oauth-go/oauth"
	"github.com/a-soliman/bookstore_users_api/domain/users"
	"github.com/a-soliman/bookstore_users_api/services"
	"github.com/a-soliman/bookstore_utils-go/rest_errors"
	"github.com/gin-gonic/gin"
)

func getUserID(userIDParam string) (int64, rest_errors.RestErr) {
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		return 0, rest_errors.NewBadRequestError("invalid user id")
	}
	return userID, nil
}

// Get finds a user by id
func Get(c *gin.Context) {
	if err := oauth.AuthenticateRequest(c.Request); err != nil {
		c.JSON(err.Status(), err)
		return
	}

	userID, idErr := getUserID(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}
	user, userErr := services.UsersService.GetUser(userID)
	if userErr != nil {
		c.JSON(userErr.Status(), userErr)
		return
	}
	// if the user is asking for their personal information, returns the private version of user data
	if oauth.GetCallerID(c.Request) == userID {
		c.JSON(http.StatusOK, user.Marshall(false))
		return
	}
	c.JSON(http.StatusOK, user.Marshall(oauth.IsPublic(c.Request)))
}

// Search finds a list of users by given search query
func Search(c *gin.Context) {
	status := strings.TrimSpace(c.Query("status"))
	fmt.Println(status)
	if status == "" {
		err := rest_errors.NewBadRequestError("missing status")
		c.JSON(err.Status(), err)
		return
	}
	result, findErr := services.UsersService.SearchUser(status)
	if findErr != nil {
		c.JSON(findErr.Status(), findErr)
		return
	}
	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-PUBLIC") == "true"))
}

// Login logs user in by email and password
func Login(c *gin.Context) {
	var request users.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		jsonErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(jsonErr.Status(), jsonErr)
		return
	}

	user, loginErr := services.UsersService.LoginUser(request)
	if loginErr != nil {
		c.JSON(loginErr.Status(), loginErr)
		return
	}
	c.JSON(http.StatusOK, user.Marshall(c.GetHeader("X-PUBLIC") == "true"))
}

// Create creates a new user
func Create(c *gin.Context) {
	var user users.User

	// read the body from request, and bind user JSON
	if err := c.ShouldBindJSON(&user); err != nil {
		jsonErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(jsonErr.Status(), jsonErr)
		return
	}

	result, saveErr := services.UsersService.CreateUser(&user)
	if saveErr != nil {
		c.JSON(saveErr.Status(), saveErr)
		return
	}
	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-PUBLIC") == "true"))
}

// Update updates a user entity
func Update(c *gin.Context) {
	userID, idErr := getUserID(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	var user users.User

	// read the body from request, and bind user JSON
	if err := c.ShouldBindJSON(&user); err != nil {
		jsonErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(jsonErr.Status(), jsonErr)
		return
	}

	user.ID = userID

	isPartial := c.Request.Method == http.MethodPatch

	result, updateErr := services.UsersService.UpdateUser(isPartial, &user)
	if updateErr != nil {
		c.JSON(updateErr.Status(), updateErr)
		return
	}
	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-PUBLIC") == "true"))
}

// Delete delete a given user
func Delete(c *gin.Context) {
	userID, idErr := getUserID(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	user := users.User{ID: userID}

	result, deleteErr := services.UsersService.DeleteUser(&user)
	if deleteErr != nil {
		c.JSON(deleteErr.Status(), deleteErr)
		return
	}
	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-PUBLIC") == "true"))
}
