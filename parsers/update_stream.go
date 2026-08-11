package parsers

import (
	"idtp/utils"
	"idtp/values"
)

type UpdateStream struct {
	DataType uint8
	Index    uint32
	Payload  []byte
}

// Returns an UpdateStream struct from an array of N bytes.
// Returns 0 if payload is missing or unable to evaluate.
func ParseUpdateStream(data []byte) (upstr UpdateStream, bytesEvaluated byte, error byte) {
	indexLength := extractIndexLength(data[0])
	dataType := extractDataType(data[0])
	dataTypeSize := utils.SizeOf(dataType)
	blockSize := 1 + indexLength + dataTypeSize

	if len(data) < int(blockSize) {
		return UpdateStream{}, 0, values.RC_MISSING_PAYLOAD
	}

	if dataTypeSize == 0 {
		return UpdateStream{}, 0, values.RC_UNKNOWN_DATA_TYPE
	}

	return UpdateStream{
		DataType: dataType,
		Index:    utils.BytesToU32(data[1 : 1+indexLength]),
		Payload:  data[1+indexLength:],
	}, blockSize, values.RC_SUCCESS
}
