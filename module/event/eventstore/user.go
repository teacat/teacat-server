package eventstore

import (
	"fmt"
	"net"
	"os"
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
		switch err.(type) {
		case *os.SyscallError, syscall.Errno, *net.OpError:
			fmt.Println("Ee")
		}
		logrus.Warningln(err)
		logrus.Warningln("Error occurred while creating an empty stream.")
	}
	return err
}
