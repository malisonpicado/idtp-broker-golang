package values

// RESPONSE CODES
const (
	RC_SUCCESS                                byte = 0x00
	RC_INVALID_INDEX                          byte = 0x01
	RC_DATA_TYPE_OVERWRITE_NOT_ALLOWED        byte = 0x02
	RC_MISSING_PAYLOAD                        byte = 0x03
	RC_UNKNOWN_METHOD                         byte = 0x04
	RC_UNSUPPORTED_PROTOCOL_VERSION           byte = 0x05
	RC_INVALID_PROTOCOL_VERSION               byte = 0x06
	RC_UNKNOWN_ENTITY_TYPE                    byte = 0x07
	RC_INVALID_KEEP_ALIVE                     byte = 0x08
	RC_CONNECTION_MUST_AUTHENTICATE           byte = 0x09
	RC_FAILED_AUTHENTICATION                  byte = 0x0A
	RC_DEVICE_MUST_DECLARE_PARAMETERS         byte = 0x0B
	RC_ONE_TIME_CONNECTION_NOT_ALLOWED        byte = 0x0C
	RC_UNKNOWN_DATA_TYPE                      byte = 0x0D
	RC_INVALID_DATA_TYPE_SET_TO_DEFAULT       byte = 0x0E
	RC_VARIABLE_OPERATION_NOT_ALLOWED         byte = 0x0F
	RC_INVALID_PARAMETER                      byte = 0x10
	RC_UNKNOWN_CONNECTION_CODE                byte = 0x11
	RC_SUCCESSFUL_CONN_WITH_PREDEFINED_PARAMS byte = 0x12
	RC_EXPANSION_LIMIT_REACHED                byte = 0x1E
)

// DATA TYPES
const (
	BOOLEAN byte = 0x00
	UINT8   byte = 0x01
	UINT16  byte = 0x02
	UINT32  byte = 0x03
	UINT64  byte = 0x04
	INT8    byte = 0x05
	INT16   byte = 0x06
	INT32   byte = 0x07 // DEFAULT
	INT64   byte = 0x08
	FLOAT32 byte = 0x09
	FLOAT64 byte = 0x0A
)

// METHODS
const (
	GET      byte = 0x00
	UPDATE   byte = 0x01
	EXPAND   byte = 0x02
	SET_TYPE byte = 0x03
)

// OPERATION MODES
const (
	OP_MODE_STRICT  byte = 0
	OP_MODE_DEFAULT byte = 1
	OP_MODE_FREE    byte = 2
)

// ENTITY TYPES
const (
	ENTITY_DEVICE byte = 0
	ENTITY_CLIENT byte = 1
)
