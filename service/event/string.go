package event

import (
	"github.com/jetbasrawi/go.geteventstore"
)

type String struct {
	Input  string
	Output string
}

func (str String) Send(option Option) error {
	event := goes.NewEvent(goes.NewUUID(), option.Stream, str, option.Meta)
	writer := option.Client.NewStreamWriter(option.Stream)

	if err := writer.Append(nil, event); err != nil {
		return err
	}

	return nil
}
