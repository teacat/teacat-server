package eventstore

import (
	"fmt"
	"net"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/TeaMeow/KitSvc/model"
	"github.com/jetbasrawi/go.geteventstore"
)

func (es *eventstore) UserCreated(u *model.User) error {
	writer := es.NewStreamWriter("user.created")

	// Create am empty stream.
	err := writer.Append(nil, goes.NewEvent("", "", u, map[string]string{}))
	if err != nil {
		switch t := err.(type) {
		case *net.OpError:
			if t.Op == "dial" {
				fmt.Println("Unknown host")
			} else if t.Op == "read" {
				fmt.Println("Connection refused")
			}
		case syscall.Errno:
			fmt.Println("Ee")
		}
		logrus.Warningln(err)
		logrus.Warningln("Error occurred while creating an empty stream.")
	}
	return err
}
