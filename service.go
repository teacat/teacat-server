package main

import (
	"errors"
	"strings"
)

var (
	// ErrEmpty 會在傳入一個空字串時被觸發。
	ErrEmpty = errors.New("字串是空的。")
)

// StringService 是基於字串的服務。
type StringService interface {
	Uppercase(string) (string, error)
	Lowercase(string) (string, error)
	Count(string) int
}

// stringService 概括了字串服務所可用的函式。
type stringService struct{}

// ServiceMiddleware 是處理 StringService 的中介層。
type ServiceMiddleware func(StringService) StringService

// Uppercase 將傳入的字串轉換為大寫。
func (stringService) Uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}

	return strings.ToUpper(s), nil
}

// Lowercase 將傳入的字串轉換為小寫。
func (stringService) Lowercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}

	return strings.ToLower(s), nil
}

// Count 計算傳入的字串長度。
func (stringService) Count(s string) int {
	return len(s)
}
