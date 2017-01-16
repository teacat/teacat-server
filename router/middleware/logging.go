package middleware

import "github.com/gin-gonic/gin"

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		//logrus.Infoln(c.Request)
		c.Next()
	}
}
