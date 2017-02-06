package server

import (
	"net/http"
	"strconv"

	"github.com/TeaMeow/KitSvc/model"
	"github.com/TeaMeow/KitSvc/module/event"
	"github.com/TeaMeow/KitSvc/module/mq"
	"github.com/TeaMeow/KitSvc/shared/auth"
	"github.com/TeaMeow/KitSvc/shared/token"
	"github.com/TeaMeow/KitSvc/store"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/olahol/melody"
)

func Created(c *gin.Context) {
	var u model.User
	if err := c.Bind(&u); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	//fmt.Println(u)
}

//
func CreateUser(c *gin.Context) {
	//time.Sleep(3 * time.Second)
	//c.AbortWithError(http.StatusInternalServerError, errors.New("Wow"))
	//return
	// Binding the data with the user struct.
	//
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
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	// Insert the user to the database.
	if err := store.CreateUser(c, &u); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	//
	//if err := mq.SendMail(c, &u); err != nil {
	//	c.AbortWithError(http.StatusInternalServerError, err)
	//	return
	//}
	go mq.SendMail(c, &u)

	go event.UserCreated(c, &u)

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

//
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
	if err := auth.Encrypt(&u.Password); err != nil {
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

//
func Login(c *gin.Context) {
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
	t, err := token.Sign(c, token.Content{ID: d.ID, Username: d.Username}, "")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": t})
}

func WebSocket(c *gin.Context) {
	w, _ := c.Get("websocket")
	ws := w.(melody.Melody)

	ws.Broadcast([]byte("Wow"))
}

func SendMail(c *gin.Context) {

}
