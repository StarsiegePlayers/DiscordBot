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

func (h *IPCHub) Register(name string, node IPCNode, ctxIn context.Context) {
	ctx, cancelFunc := context.WithCancel(ctxIn)
	h.Nodes[name] = IPCNodeContext{
		IPCNode: node,
		Cancel:  cancelFunc,
	}

	go h.listen(ctx, node.RX)
}

func (h *IPCHub) Process(mainHandler IPCHandler) {
	for {
		select {
		case m := <-h.MainHub.RX:
			go mainHandler(m)
			h.BroadcastMessage(m)

		case <-h.Context.Done():
			log.Println("[IPC Hub]: shutting down")

			for k, v := range h.Nodes {
				log.Printf("[IPC Hub] canceling %s\n", k)
				v.Cancel()
			}

			return
		}
	}
}

func (h *IPCHub) listen(ctx context.Context, ipxRX chan IPCMessage) {
	for {
		select {
		case m := <-ipxRX:
			h.MainHub.RX <- m
		case <-ctx.Done():
			return
		}
	}
}

func (h *IPCHub) BroadcastMessage(m IPCMessage) {
	for _, v := range h.Nodes {
		v.TX <- m
	}

	log.Printf("[IPC-Broadcast]: from [%s] {%s}", m.Sender, m)
}
