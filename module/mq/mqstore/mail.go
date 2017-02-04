package mqstore

import (
	"encoding/json"

	"github.com/TeaMeow/KitSvc/model"
)

func (m *mqstore) SendMail(user *model.User) error {
	return m.send("send_mail", user)
}

func (m *mqstore) send(topic string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return m.Publish(topic, b)
}
