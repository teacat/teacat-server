package main

import "strings"

// Model function represents the business logic and processes with the data.
//
// Create the model functions with the following format:
//     func (m Model)...

// ToUpper converts the string to uppercase.
func (m Model) ToUpper(s string) (string, error) {
	if s == "" {
		return "", Err{
			Message: ErrEmpty,
		}
	}

	return strings.ToUpper(s), nil
}

// Count counts the length of the string.
func (m Model) Count(s string) int {
	return len(s)
}
