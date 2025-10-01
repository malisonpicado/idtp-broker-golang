package values

type StatusCode uint8

const (
	SUCCESS                          StatusCode = 0x00
	INVALID_INDEX                    StatusCode = 0x01
	DATA_TYPE_OVERWRITE_NOT_ALLOWED  StatusCode = 0x02
	MISSING_PAYLOAD                  StatusCode = 0x03
	UNKNOWN_METHOD                   StatusCode = 0x04
	UNSUPPORTED_PROTOCOL_VERSION     StatusCode = 0x05
	INVALID_PROTOCOL_VERSION         StatusCode = 0x06
	UNKNOWN_ENTITY_TYPE              StatusCode = 0x07
	INVALID_KEEP_ALIVE               StatusCode = 0x08
	CONNECTION_MUST_AUTHENTICATE     StatusCode = 0x09
	FAILED_AUTHENTICATION            StatusCode = 0x0A
	DEVICE_MUST_DECLARE_PARAMETERS   StatusCode = 0x0B
	ONTE_TIME_CONNECTION_NOT_ALLOWED StatusCode = 0x0C
	UNKNOWN_DATA_TYPE                StatusCode = 0x0D
	INVALID_DATA_TYPE_SET_TO_DEFAULT StatusCode = 0x0E
	EXPANSION_LIMIT_REACHED          StatusCode = 0x1E
)

type DataType uint8

const (
	BOOLEAN DataType = 0x00
	UINT8   DataType = 0x01
	UINT16  DataType = 0x02
	UINT32  DataType = 0x03
	UINT64  DataType = 0x04
	INT8    DataType = 0x05
	INT16   DataType = 0x06
	INT32   DataType = 0x07 // DEFAULT
	INT64   DataType = 0x08
	FLOAT32 DataType = 0x09
	FLOAT64 DataType = 0x0A
)

type MethodCode uint8

const (
	GET      MethodCode = 0x00
	UPDATE   MethodCode = 0x01
	EXPAND   MethodCode = 0x02
	SET_TYPE MethodCode = 0x03
)

type OperationMode uint8

const (
	OP_MODE_STRICT  OperationMode = 0
	OP_MODE_DEFAULT OperationMode = 1
	OP_MODE_FREE    OperationMode = 2
)

type EntityType uint8

const (
	ENTITY_DEVICE EntityType = 0
	ENTITY_CLIENT EntityType = 1
)

type Configuration struct {
	Key             string
	ProtocolVersion uint8
	OperationMode   OperationMode
	StorageLimit    uint32
}

func SizeOf(t DataType) uint8 {
	switch t {
	case BOOLEAN, UINT8, INT8:
		return 1
	case UINT16, INT16:
		return 2
	case UINT32, INT32, FLOAT32:
		return 4
	case UINT64, INT64, FLOAT64:
		return 8
	}

	return 0
}
