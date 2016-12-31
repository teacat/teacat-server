package service

import (
	"errors"
	"fmt"
	"net/http"

	nsq "github.com/bitly/go-nsq"
)

var (
	// ErrEmpty 會在傳入一個空字串時被觸發。
	ErrEmpty = ErrInfo{
		Text:   errors.New("The string is empty."),
		Status: http.StatusBadRequest,
		Code:   "str_empty",
	}
)

// StringService 是基於字串的服務。
type Service interface {
	Uppercase(string) (string, error)
	Count(string) int
	Test(*nsq.Message)
}

// Uppercase 將傳入的字串轉換為大寫。
func (c Concrete) Uppercase(s string) (string, error) {

	//c.Message.Publish("new_user", []byte("test"))

	res, err := c.Model.ToUpper(s)
	if err != nil {
		return "", err
	}

	return res, nil
}

// Count 計算傳入的字串長度。
func (c Concrete) Count(s string) int {
	return c.Model.Count(s)
}

func (Concrete) Test(msg *nsq.Message) {
	fmt.Println(msg)
}
