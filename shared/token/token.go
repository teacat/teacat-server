package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/TeaMeow/KitSvc/model"
	"github.com/TeaMeow/KitSvc/router/middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	ErrMissingHeader = errors.New("The length of the `Authorization` header is zero.")
)

type Content struct {
	ID       int
	Username string
}

func secretFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		// Make sure the `alg` is what we except.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(secret), nil
	}
}

func Parse(tokenString string, secret string) (*Content, error) {
	content := &Content{}

	// Parse the token.
	token, err := jwt.Parse(tokenString, secretFunc(secret))

	// Parse error.
	if err != nil {
		return content, err

		// Read the token if it's valid.
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		content.ID = int(claims["id"].(float64))
		content.Username = claims["username"].(string)
		return content, nil

		// Other errors.
	} else {
		return content, err
	}
}

func ParseRequest(c *gin.Context) (*Content, error) {

	// Get the jwt from the `Authorization` header.
	header := c.Request.Header.Get("Authorization")

	// Load the jwt secret from the gin config
	config, _ := c.Get(middleware.ConfigKey)
	secret := config.(*model.Config).JWTSecret

	if len(header) == 0 {
		return &Content{}, ErrMissingHeader
	}

	var t string
	// Parse the header to get the token part.
	fmt.Sscanf(header, "Bearer %s", &t)
	return Parse(t, secret)
}

func Sign(ctx *gin.Context, c Content, secret string) (tokenString string, err error) {

	// Load the jwt secret from the Gin config if the secret isn't specified.
	if secret == "" {
		config, _ := ctx.Get(middleware.ConfigKey)
		secret = config.(*model.Config).JWTSecret
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       c.ID,
		"username": c.Username,
		"nbf":      time.Now().Unix(),
		"iat":      time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret.
	tokenString, err = token.SignedString([]byte(secret))

	return
}
