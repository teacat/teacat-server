package service

import (
	"net/http"
	"strconv"
	"time"

	"github.com/TeaMeow/KitSvc/errno"
	"github.com/TeaMeow/KitSvc/model"
	"github.com/TeaMeow/KitSvc/module/event"
	"github.com/TeaMeow/KitSvc/module/mq"
	"github.com/TeaMeow/KitSvc/shared/auth"
	"github.com/TeaMeow/KitSvc/shared/token"
	"github.com/TeaMeow/KitSvc/shared/wsutil"
	"github.com/TeaMeow/KitSvc/store"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/olahol/melody"
)

// CreateUser creates a new user account.
func CreateUser(c *gin.Context) {
	var u model.User
	if err := c.Bind(&u); err != nil {
		errno.Abort(errno.ErrBind, err, c)
		return
	}
	// Validate the data.
	if err := u.Validate(); err != nil {
		errno.Abort(errno.ErrValidation, err, c)
		return
	}
	// Encrypt the user password.
	if err := u.Encrypt(); err != nil {
		errno.Abort(errno.ErrEncrypt, err, c)
		return
	}
	// Insert the user to the database.
	if err := store.CreateUser(c, &u); err != nil {
		errno.Abort(errno.ErrDatabase, err, c)
		return
	}
	// Publish the `send_mail` message to the message queue.
	mq.Publish(c, mq.M{
		Topic: mq.MsgSendMail,
		Data:  &u,
	})
	// Send the `user_created` event to Event Store.
	event.Send(c, event.E{
		Stream: event.EvtUserCreated,
		Data:   &u,
	})
	// Show the user information.
	c.JSON(http.StatusOK, u)
}

// GetUser gets an user by the user identifier.
func GetUser(c *gin.Context) {
	// Get the `username` from the url parameter.
	username := c.Param("username")
	// Get the user by the `username` from the database.
	u, err := store.GetUser(c, username)
	if err != nil {
		errno.Abort(errno.ErrUserNotFound, err, c)
		return
	}

	c.JSON(http.StatusOK, u)
}

// DeleteUser deletes the user by the user identifier.
func DeleteUser(c *gin.Context) {
	// Get the user id from the url parameter.
	userID, _ := strconv.Atoi(c.Param("id"))
	// Delete the user in the database.
	if err := store.DeleteUser(c, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			errno.Abort(errno.ErrUserNotFound, err, c)
		} else {
			errno.Abort(errno.ErrDatabase, err, c)
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
		errno.Abort(errno.ErrBind, err, c)
		return
	}

	// We update the record based on the user id.
	u.ID = userID
	// Validate the data.
	if err := u.Validate(); err != nil {
		errno.Abort(errno.ErrValidation, err, c)
		return
	}
	// Parse the json web token.
	if _, err := token.ParseRequest(c); err != nil {
		errno.Abort(errno.ErrTokenInvalid, err, c)
		return
	}
	// Encrypt the user password.
	if err := u.Encrypt(); err != nil {
		errno.Abort(errno.ErrEncrypt, err, c)
		return
	}
	// Update the user in the database.
	if err := store.UpdateUser(c, &u); err != nil {
		if err == gorm.ErrRecordNotFound {
			errno.Abort(errno.ErrUserNotFound, err, c)
		} else {
			errno.Abort(errno.ErrDatabase, err, c)
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
		errno.Abort(errno.ErrBind, err, c)
		return
	}
	// Get the user information by the login username.
	d, err := store.GetUser(c, u.Username)
	if err != nil {
		errno.Abort(errno.ErrUserNotFound, err, c)
		return
	}
	// Compare the login password with the user password.
	if err := auth.Compare(d.Password, u.Password); err != nil {
		errno.Abort(errno.ErrPasswordIncorrect, err, c)
		return
	}
	// Sign the json web token.
	t, err := token.Sign(c, token.Context{ID: d.ID, Username: d.Username}, "")
	if err != nil {
		errno.Abort(errno.ErrToken, err, c)
		return
	}

	c.JSON(http.StatusOK, model.Token{Token: t})
}

// WatchUser watches the user changes, and broadcast when there's a new user.
func WatchUser(c *gin.Context) {
	// Get the WebSocket from the Gin context.
	ws := wsutil.Get(c)

	ws.HandleConnect(func(s *melody.Session) {
		go func() {
			// Get the identifier of the last user.
			l, _ := store.GetLastUser(c)
			// Keep watching the user changes.
			for {
				// Slow and steady, otherwise the database will boomed.
				<-time.After(100 * time.Millisecond)
				// Get the new user who is registered after the specified user.
				u, err := store.GetUserAfter(c, l.ID)
				if err != nil {
					continue
				}
				// Boardcast the lastest username.
				ws.Broadcast([]byte(u.Username))

				l.ID++
			}
		}()
	})
}

// SendMail sends the mail to the new user's inbox.
func SendMail(c *gin.Context) {
	// Blah, blah blah ...
}

// UserCreated handles the `user-created` event.
func UserCreated(c *gin.Context) {
	var u model.User
	if err := c.Bind(&u); err != nil {
		errno.Abort(errno.ErrBind, err, c)
		return
	}
}
