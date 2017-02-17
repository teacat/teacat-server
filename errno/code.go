package errno

import "net/http"

var (
	ErrBind = &Err{
		Code:       "BIND_ERR",
		Message:    "Error occurred while binding the request body to the struct.",
		StatusCode: http.StatusBadRequest}
	ErrValidation = &Err{
		Code:       "VALIDATION_ERR",
		Message:    "Validation failed.",
		StatusCode: http.StatusBadRequest}
	ErrEncrypt = &Err{
		Code:       "ENCRYPT_ERR",
		Message:    "Error occurred while encrypting the user password.",
		StatusCode: http.StatusInternalServerError}
	ErrDatabase = &Err{
		Code:       "DB_ERR",
		Message:    "Database error.",
		StatusCode: http.StatusInternalServerError}
	ErrUserNotFound = &Err{
		Code:       "USER_NOT_FOUND",
		Message:    "The user was not found.",
		StatusCode: http.StatusNotFound}
	ErrTokenInvalid = &Err{
		Code:       "TOKEN_INVALID",
		Message:    "The token was invalid.",
		StatusCode: http.StatusForbidden}
	ErrPasswordIncorrect = &Err{
		Code:       "PASSWORD_INCORRECT",
		Message:    "The password was incorrect.",
		StatusCode: http.StatusForbidden}
	ErrToken = &Err{
		Code:       "TOKEN_ERR",
		Message:    "Error occurred while signing the JSON web token.",
		StatusCode: http.StatusBadRequest}
)
