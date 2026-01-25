package values

type Request struct {
	Method   byte
	DataType byte
	Index    uint32
	Payload  []byte
}

type DeviceParameter struct {
	Method byte
	Index  uint32
}

type ConnectionRequest struct {
	ProtocolVersion byte
	EntityType      byte
	KeepAlive       uint16
	HasAuth         bool
	UserKey         string
	Parameters      []DeviceParameter
}

type Configuration struct {
	Key             string
	ProtocolVersion byte
	OperationMode   byte
	StorageLimit    uint32
}

type Entity struct {
	EntityType       byte
	ProcessAsStrict  bool
	UpdateParams     []uint32
	DependencyParams []uint32
}
