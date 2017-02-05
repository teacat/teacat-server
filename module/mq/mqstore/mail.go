package mqstore

import "github.com/TeaMeow/KitSvc/model"

func (m *mqstore) SendMail(user *model.User) error {
	m.send("send_mail", user)
	return nil
}
