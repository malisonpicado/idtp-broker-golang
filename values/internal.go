package values

type ProcessAs struct {
	Type      ProcessType
	PreRcCode byte
}

type ProcessType byte

const (
	PT_ERROR             ProcessType = 1
	PT_REQUEST           ProcessType = 2
	PT_CREATE_CONNECTION ProcessType = 3
	PT_PING              ProcessType = 4
	PT_DISCONNECTION     ProcessType = 5
)
