package middleware

import (
	"fmt"

	"github.com/codegangsta/cli"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWT(c *cli.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Get the token from the header.
		auth := ctx.Request.Header.Get("Authorization")
		if len(auth) < 7 {
			ctx.JSON(400, gin.H{
				"status":  "error",
				"code":    "token_empty",
				"message": "The token was not found in the request header.",
				"payload": nil,
			})
			return
		}

		// Remove the `Bearer ` leading, leave the token part only.
		tokenString := auth[7:len(auth)]

		// Parse the token, and make sure the `alg` is correct.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Make sure the `alg`` is what we except.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(c.String("jwt-secret")), nil
		})
		if err != nil {
			ctx.JSON(400, gin.H{
				"status":  "error",
				"code":    "token_error",
				"message": err.Error(),
				"payload": nil,
			})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println(claims["foo"], claims["nbf"])
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				ctx.JSON(400, gin.H{
					"status":  "error",
					"code":    "token_incorrect",
					"message": "The token format is incorrect.",
					"payload": nil,
				})
				return
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				ctx.JSON(400, gin.H{
					"status":  "error",
					"code":    "token_expired",
					"message": "The token has been expired or not active yet.",
					"payload": nil,
				})
				return
			} else {
				ctx.JSON(400, gin.H{
					"status":  "error",
					"code":    "token_error",
					"message": err.Error(),
					"payload": nil,
				})
				return
			}
		} else {
			ctx.JSON(400, gin.H{
				"status":  "error",
				"code":    "token_error",
				"message": err.Error(),
				"payload": nil,
			})
			return
		}

		//c.Set()
		ctx.Next()
	}
}
