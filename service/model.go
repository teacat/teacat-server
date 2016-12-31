package service

import "strings"

func (m Model) ToUpper(s string) (string, error) {
	if s == "" {
		return "", Err{
			Message: ErrEmpty,
		}
	}

	return strings.ToUpper(s), nil
}

func (m Model) Count(s string) int {
	return len(s)
}
