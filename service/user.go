package service

import (
	"net/http"
	"strconv"

	"github.com/TeaMeow/KitSvc/model"
	"github.com/TeaMeow/KitSvc/module/event"
	"github.com/TeaMeow/KitSvc/module/logger"
	"github.com/TeaMeow/KitSvc/module/mq"
	"github.com/TeaMeow/KitSvc/shared/auth"
	"github.com/TeaMeow/KitSvc/shared/token"
	"github.com/TeaMeow/KitSvc/shared/wsutil"
	"github.com/TeaMeow/KitSvc/store"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// CreateUser creates a new user account.
func CreateUser(c *gin.Context) {
	var u model.User

	if err := c.Bind(&u); err != nil {
		c.
			AbortWithError(http.StatusBadRequest, err).
			SetMeta(logger.Meta("BIND_ERROR"))
		return
	}
	// Validate the data.
	if err := u.Validate(); err != nil {
		c.
			String(http.StatusBadRequest, err.Error())
		return
	}
	// Encrypt the user password.
	if err := u.Encrypt(); err != nil {
		c.
			AbortWithError(http.StatusInternalServerError, err)
		return
	}
	// Insert the user to the database.
	if err := store.CreateUser(c, &u); err != nil {
		c.
			AbortWithError(http.StatusInternalServerError, err)
		return
	}
	// Publish the `send_mail` message to the message queue.
	go mq.SendMail(c, &u)
	// Send a `user_created` event to Event Store.
	go event.UserCreated(c, &u)

	// Show the user information.
	c.JSON(http.StatusOK, u)
}

// GetUser gets an user by the user identifier.
func GetUser(c *gin.Context) {
	// Get the `username` from the url parameter.
	username := c.Param("username")

	// Get the user by the `username` from the database.
	if u, err := store.GetUser(c, username); err != nil {
		c.String(http.StatusNotFound, "The user was not found.")
	} else {
		c.JSON(http.StatusOK, u)
	}
}

// DeleteUser deletes the user by the user identifier.
func DeleteUser(c *gin.Context) {
	// Get the user id from the url parameter.
	userID, _ := strconv.Atoi(c.Param("id"))
	// Delete the user in the database.
	if err := store.DeleteUser(c, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.String(http.StatusNotFound, "The user doesn't exist.")
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	c.String(http.StatusOK, "The user has been deleted successfully.")
}

// UpdateUser updates an user account information.
func UpdateUser(c *gin.Context) {
	// Get the user id from the url parameter.
	userID, _ := strconv.Atoi(c.Param("id"))

	// Binding the user data.
	var u model.User
	if err := c.Bind(&u); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// We update the record based on the user id.
	u.ID = userID

	// Validate the data.
	if err := u.Validate(); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	// Parse the json web token.
	if _, err := token.ParseRequest(c); err != nil {
		c.String(http.StatusForbidden, "The token was incorrect.")
		return
	}
	// Encrypt the user password.
	if err := u.Encrypt(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	// Update the user in the database.
	if err := store.UpdateUser(c, &u); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.String(http.StatusNotFound, "The user doesn't exist.")
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	c.JSON(http.StatusOK, u)
}

// PostToken generates the authentication token
// if the password was matched with the specified account.
func PostToken(c *gin.Context) {
	// Binding the data with the user struct.
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
	t, err := token.Sign(c, token.Context{ID: d.ID, Username: d.Username}, "")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, model.Token{Token: t})
}

// WatchUser watches the user changes, and broadcast when there's a new user.
func WatchUser(c *gin.Context) {
	ws := wsutil.Get(c)

	ws.Broadcast([]byte("Wow"))
}

// SendMail sends the mail to the new user's inbox.
func SendMail(c *gin.Context) {
	// Blah, blah blah ...
}

// UserCreated handles the `user-created` event.
func UserCreated(c *gin.Context) {
	var u model.User
	if err := c.Bind(&u); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
}
