package errno

import "net/http"

var (
	errs = map[string]*Err{
		"BIND_ERR":       {Message: "Error occurred while binding the request body to the struct.", StatusCode: http.StatusBadRequest},
		"VALIDATION_ERR": {Message: "Validation failed.", StatusCode: http.StatusBadRequest},
		"ENCRYPT_ERR":    {Message: "Error occurred while encrypting the user password.", StatusCode: http.StatusInternalServerError},
		"DB_ERR":         {Message: "Database error.", StatusCode: http.StatusInternalServerError},
		"USER_NOT_FOUND": {Message: "The user was not found.", StatusCode: http.StatusNotFound},
        "TOKEN_INVALID": {Message: "The token was invalid.", StatusCode: http.StatusForbidden},
        "PASSWORD_INCORRECT": {Message: "The password was incorrect.", StatusCode: http.StatusForbidden},
        "TOKEN_ERR": {Message: "Error occurred while signing the JSON web token.", StatusCode: http.StatusBadRequest},
	}
)
