package plugin

import (
	"fmt"
)

type IPCHandler func(IPCMessage)

type IPCMessage struct {
	Sender      string
	Destination string
	Command     IPCCommand
	Message     string
}

func (m IPCMessage) String() string {
	return fmt.Sprintf("{%s} [%s] %s", IPCCommandStrings[m.Command], m.Destination, m.Message)
}
