package messaging

import "github.com/TeaMeow/KitSvc/service"

func (conc Concrete) SetHandlers(svc service.Service) {

	/*testHandler := func(message *nsq.Message) error {
		log.Printf("Got a message: %v", message)
		return nil
	}*/

	//conc.Handle("new_user", "string", svc.CreateString)
}
