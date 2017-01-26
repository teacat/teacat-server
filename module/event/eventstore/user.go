package eventstore

import "github.com/TeaMeow/KitSvc/model"

func (es *eventstore) UserCreated(u *model.User) error {
	es.send("user.created", u, map[string]string{})

	return nil
}
