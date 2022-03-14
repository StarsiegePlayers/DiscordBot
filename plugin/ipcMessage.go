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
	return fmt.Sprintf("[IPC]: [%s] %s", m.Destination, m.Message)
}
