package server

import (
	"net/http"
	"strconv"

	"github.com/TeaMeow/KitSvc/model"
	"github.com/TeaMeow/KitSvc/shared/auth"
	"github.com/TeaMeow/KitSvc/store"
	"github.com/gin-gonic/gin"
)

//
func CreateUser(c *gin.Context) {
	var u model.User
	if err := c.Bind(&u); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
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
	}

	if err := auth.Compare(d.Password, u.Password); err != nil {
		c.String(http.StatusForbidden, "The username or the password was incorrect.")
	} else {
		c.String(200, strconv.Itoa(d.ID))
	}
}
