package event

import "github.com/jetbasrawi/go.geteventstore"

type Option struct {
	Client *goes.Client
	Stream string
	Meta   map[string]string
}
