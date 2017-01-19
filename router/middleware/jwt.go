package middleware

import (
	"fmt"

	"github.com/codegangsta/cli"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const jwtKey = "jwt"

func JWT(c *cli.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Get the token from the header.
		auth := ctx.Request.Header.Get("Authorization")
		if len(auth) < 7 {
			ctx.JSON(400, gin.H{
				"status":  "error",
				"code":    "token_error",
				"message": "The token was not found or the length was incorrect.",
				"payload": nil,
			})
			return
		}

		// Remove the `Bearer ` leading, leave the token part only.
		tokenString := auth[7:len(auth)]
		// Parse the token, and make sure the `alg` is correct.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Make sure the `alg` is what we except.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			// Return the secret if the token has a valid alg.
			return []byte(c.String("jwt-secret")), nil
		})

		// Parse error.
		if err != nil {
			ctx.JSON(400, gin.H{
				"status":  "error",
				"code":    "token_error",
				"message": err.Error(),
				"payload": nil,
			})

			// Read the token if it's valid.
		} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Put the JWT in Gin context, so we can get the token payload from the hanlders.
			ctx.Set(jwtKey, claims)
			ctx.Next()

			// Token validation error.
		} else if ve, ok := err.(*jwt.ValidationError); ok {

			// Incorrect token format.
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				ctx.JSON(400, gin.H{
					"status":  "error",
					"code":    "token_error",
					"message": "The token format is incorrect.",
					"payload": nil,
				})

				// Expired or not active yet.
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				ctx.JSON(400, gin.H{
					"status":  "error",
					"code":    "token_expired",
					"message": "The token has been expired or not active yet.",
					"payload": nil,
				})

				// Other errors.
			} else {
				ctx.JSON(400, gin.H{
					"status":  "error",
					"code":    "token_error",
					"message": err.Error(),
					"payload": nil,
				})
			}

			// Other errors.
		} else {
			ctx.JSON(400, gin.H{
				"status":  "error",
				"code":    "token_error",
				"message": err.Error(),
				"payload": nil,
			})
		}

		return
	}
}
