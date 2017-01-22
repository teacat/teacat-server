package server

import (
	"net/http"
	"strconv"

	"github.com/TeaMeow/KitSvc/model"
	"github.com/TeaMeow/KitSvc/shared/auth"
	"github.com/TeaMeow/KitSvc/shared/token"
	"github.com/TeaMeow/KitSvc/store"
	"github.com/gin-gonic/gin"
)

//
func CreateUser(c *gin.Context) {
	// Binding the data with the user struct.
	var u model.User
	if err := c.Bind(&u); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// Validate the data.
	if err := u.Validate(); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	// Encrypt the user password.
	if err := auth.Encrypt(&u.Password); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// Insert the user to the database.
	if err := store.CreateUser(c, &u); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// Show the user information.
	c.JSON(http.StatusOK, u)
}

//
func GetUser(c *gin.Context) {
	// Get the `username` from the url parameter.
	username := c.Param("username")

	// Get the user by the `username`` from the database.
	if u, err := store.GetUser(c, username); err != nil {
		c.String(http.StatusNotFound, "The user was not found.")
	} else {
		c.JSON(http.StatusOK, u)
	}
}

//
func DeleteUser(c *gin.Context) {

}

//
func UpdateUser(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))

	t, err := token.ParseRequest(c)
	if err != nil {
		c.String(http.StatusForbidden, "The token was incorrect.")
		return
	} else {
		c.JSON(http.StatusOK, t)
		return
	}

	return

	if store.Can(c, &model.Permission{
		Action:     model.PERM_EDIT,
		ResourceID: userID,
		UserID:     userID,
	}) {
		c.String(http.StatusOK, "Okay: %d", userID)
		return
	}

	c.String(http.StatusForbidden, "What the fuck? %d", userID)
}

//
func Login(c *gin.Context) {
	var u model.User
	if err := c.Bind(&u); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Get the user information by the login username.
	d, err := store.GetUser(c, u.Username)
	if err != nil {
		c.String(http.StatusNotFound, "The user doesn't exist.")
		return
	}

	// Compare the login password with the user password.
	if err := auth.Compare(d.Password, u.Password); err != nil {
		c.String(http.StatusForbidden, "The password was incorrect.")
		return
	}

	// Sign the json web token.
	t, err := token.Sign(c, token.Content{ID: d.ID, Username: d.Username}, "")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.String(http.StatusOK, t)
}
