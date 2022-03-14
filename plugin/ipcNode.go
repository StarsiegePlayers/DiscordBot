package plugin

type IPCNode struct {
	TX chan IPCMessage
	RX chan IPCMessage
}
