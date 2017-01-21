package server

import (
	"net/http"
	"strconv"

	"github.com/TeaMeow/KitSvc/client"
	"github.com/TeaMeow/KitSvc/model"
	"github.com/TeaMeow/KitSvc/protobuf"
	"github.com/TeaMeow/KitSvc/shared/auth"
	"github.com/TeaMeow/KitSvc/shared/token"
	"github.com/TeaMeow/KitSvc/store"
	"github.com/gin-gonic/gin"
)

//
func CreateUser(c *gin.Context) {

	var t protobuf.CreateUserRequest
	if err := c.Bind(&t); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	u := model.User{
		Username: t.Username,
		Password: t.Password,
	}

	u.Password, _ = auth.Encrypt(u.Password)

	if err := store.CreateUser(c, &u); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.String(200, strconv.Itoa(u.ID))
}

//
func GetUser(c *gin.Context) {
	username := c.Param("username")

	cli := client.NewClient("http://localhost:8080")
	cli.PostUser(&model.User{
		Username: "Wow",
		Password: "wowowowowo",
	})

	if user, err := store.GetUser(c, username); err != nil {
		c.String(http.StatusNotFound, "The user was not found.")
	} else {
		c.JSON(http.StatusOK, user)
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
		c.String(400, "The token was incorrect.")
	} else {
		c.JSON(200, t)
	}

	return

	if store.Can(c, &model.Permission{
		Action:     model.PERM_EDIT,
		ResourceID: userID,
		UserID:     userID,
	}) {
		c.String(200, "Okay: %d", userID)
	}

	c.String(403, "What the fuck? %d", userID)
}

//
func Login(c *gin.Context) {
	var u model.User
	if err := c.Bind(&u); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	d, err := store.GetUser(c, u.Username)
	if err != nil {
		c.String(http.StatusNotFound, "The user doesn't exist.")
		return
	}

	if err := auth.Compare(d.Password, u.Password); err != nil {
		c.String(http.StatusForbidden, "The username or the password was incorrect.")
		return
	}

	t, err := token.Sign(c, token.Content{ID: d.ID, Username: d.Username}, "")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.String(200, t)
	return
}
