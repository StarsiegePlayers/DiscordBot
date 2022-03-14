package plugin

import (
	"context"
	"log"
)

type IPCNodeContext struct {
	IPCNode
	Cancel context.CancelFunc
}

type IPCHub struct {
	Nodes   map[string]IPCNodeContext
	MainHub IPCNode
	Log     log.Logger
	Context context.Context
}

func NewIPCHub(ctx context.Context) *IPCHub {
	return &IPCHub{
		Nodes: make(map[string]IPCNodeContext),
		MainHub: IPCNode{
			RX: make(chan IPCMessage, MIN_BUFFER),
		},
		Context: ctx,
	}
}

func (h *IPCHub) Register(name string, node IPCNode) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	h.Nodes[name] = IPCNodeContext{
		IPCNode: node,
		Cancel:  cancelFunc,
	}

	go h.listen(ctx, node.RX)
}

func (h *IPCHub) Process(mainHandler IPCHandler) {
	select {
	case m := <-h.MainHub.RX:
		go mainHandler(m)
		h.BroadcastMessage(m)

	case <-h.Context.Done():
		log.Println("[IPC Hub]: shutting down")

		for _, v := range h.Nodes {
			v.Cancel()
		}
	}
}

func (h *IPCHub) listen(ctx context.Context, ipxRX chan IPCMessage) {
	select {
	case m := <-ipxRX:
		h.MainHub.RX <- m
	case <-ctx.Done():
		return
	}
}

func (h *IPCHub) BroadcastMessage(m IPCMessage) {
	h.Log.Println("[IPC Broadcast]:", m)

	for _, v := range h.Nodes {
		v.TX <- m
	}
}
