package plugin

const (
	BROADCAST  = "BROADCAST"
	HUB        = "HUB"
	CONFIG     = "CONFIG"
	MIN_BUFFER = 5
)

type IPCCommand int

const (
	NULL IPCCommand = iota
	LOG
	REGISTER
	MESSAGE
)

var IPCCommandStrings = map[IPCCommand]string{
	NULL:     "",
	REGISTER: "REGISTER",
	LOG:      "LOG",
	MESSAGE:  "MESSAGE",
}
