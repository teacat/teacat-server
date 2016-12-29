package messaging

import "github.com/TeaMeow/KitSvc/service"

func (c Concrete) SetHandlers(svc service.Service) {
	c.Handle("new_user", "string", svc.Test)
}
