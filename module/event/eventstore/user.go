package eventstore

import (
	"github.com/TeaMeow/KitSvc/model"
	"github.com/TeaMeow/KitSvc/version"
)

// UserCreated handles the `user_created` event.
func (es *eventstore) UserCreated(u *model.User) error {
	es.send("user_created", u, map[string]string{"node": version.Version})

	return nil
}
