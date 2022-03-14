package plugin

const (
	BROADCAST  = "BROADCAST"
	HUB        = "HUB"
	MIN_BUFFER = 5
)

type IPCCommand int

const (
	REGISTER IPCCommand = iota
	MESSAGE
)
