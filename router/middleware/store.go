package middleware

import (
	"github.com/TeaMeow/KitSvc/store"
	"github.com/TeaMeow/KitSvc/store/datastore"
	"github.com/gin-gonic/gin"
)

func Store() gin.HandlerFunc {
	v := datastore.Open()

	return func(c *gin.Context) {
		store.ToContext(c, v)
		c.Next()
	}
}
