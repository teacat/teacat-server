package eventstore

import (
	"fmt"
	"time"

	"github.com/TeaMeow/KitSvc/model"
	"github.com/jetbasrawi/go.geteventstore"
)

func (es *eventstore) Send(stream string, data interface{}, meta interface{}) error {
	//
	writer := es.NewStreamWriter(stream)

	//
	//maxRetry := 5
	//retried := 0

	for {
		if es.isConnected {
			err := writer.Append(nil, goes.NewEvent("", "", data, meta))
			return err
		}

		fmt.Println("Waiting for connection.")
		<-time.After(time.Second * 1a)

		/*err := writer.Append(nil, goes.NewEvent("", "", data, meta))
		if err == nil {
			return nil
		}

		//
		if retried < maxRetry {
			logrus.Warningln(err)
			logrus.Warningf("Error occurred while creating the `%s` event, retry after 1 second. (%d/%d)", stream, retried, maxRetry)

			retried++

			<-time.After(time.Second * time.Duration(retried))

			//
		} else {
			logrus.Errorf("Cannot create the `%s` event after retried %d times.", stream, retried)
			return err
		}*/
	}
}

func (es *eventstore) UserCreated(u *model.User) error {
	err := es.Send("user.created", u, map[string]string{})

	return err
}
