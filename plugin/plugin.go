package plugin

import (
	"context"
	"fmt"
	"log"
	"strings"
)

type Base struct {
	CTX     context.Context
	IPCNode IPCNode
	Name    string

	commands map[string]IPCHandler
}

func (b *Base) Init(rxIn chan IPCMessage) chan IPCMessage {
	b.IPCNode.RX = rxIn
	b.IPCNode.TX = make(chan IPCMessage, MIN_BUFFER)

	return b.IPCNode.TX
}

func (b *Base) Attach(ctx context.Context) {
	select {
	case m := <-b.IPCNode.RX:
		go b.ProcessEvent(m)

	case <-ctx.Done():
		log.Printf("[%s] shutting down...", b.Name)
	}
}

func (b *Base) ProcessEvent(m IPCMessage) {
	s := strings.Fields(m.Message)
	command := s[0]

	if f, ok := b.commands[command]; ok {
		f(m)
	}
}

func (b *Base) SendMessage(m IPCMessage) {
	m.Sender = b.Name
	b.IPCNode.TX <- m
}

func (b *Base) Logf(format string, args ...interface{}) {
	fmat := fmt.Sprintf("[%s]: %s", b.Name, format)
	log.Printf(fmat+"\n", args...)

	b.SendMessage(IPCMessage{
		Destination: "LOG",
		Message:     fmt.Sprintf(fmat, args...),
	})
}

func (b *Base) RegisterCommand(name string, f IPCHandler) {
	b.commands[name] = f
	b.IPCNode.TX <- IPCMessage{
		Sender:      b.Name,
		Destination: BROADCAST,
		Command:     REGISTER,
		Message:     name,
	}
}

func (b *Base) Unload() int {
	return 0
}
