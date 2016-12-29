package service

import (
	"errors"
	"strings"

	nsq "github.com/bitly/go-nsq"
)

var (
	// ErrEmpty 會在傳入一個空字串時被觸發。
	ErrEmpty = errors.New("字串是空的。")
)

// StringService 是基於字串的服務。
type Service interface {
	Uppercase(string) (string, error)
	Lowercase(string) (string, error)
	Count(string) int
}

// stringService 概括了字串服務所可用的函式。
type Concrete struct {
	Message *nsq.Producer
}

// ServiceMiddleware 是處理 StringService 的中介層。
type Middleware func(Service) Service

// Uppercase 將傳入的字串轉換為大寫。
func (Concrete) Uppercase(s string) (string, error) {

	//c.Message.Publish("write_test", []byte("test"))

	if s == "" {
		return "", ErrEmpty
	}

	return strings.ToUpper(s), nil
}

// Lowercase 將傳入的字串轉換為小寫。
func (Concrete) Lowercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}

	return strings.ToLower(s), nil
}

// Count 計算傳入的字串長度。
func (Concrete) Count(s string) int {
	return len(s)
}

//
/*func (Concrete) Test(msg *nsq.Message) {
	string(msg)
}*/
