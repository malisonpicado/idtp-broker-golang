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

// Converts an index to byte format. Returns the index length
// code and the index in bytes based on the length code value.
// 0x00 -> 1 byte, ..., 0x03 -> 4 bytes
func CompactIndex(index uint32) (ilen byte, ibyte []byte) {
	i := utils.U32ToBytes(index)

	if index < 256 {
		return 0x00, i[3:4]
	}

	if index < 65_536 {
		return 0x01, i[2:4]
	}

	if index < 16_777_216 {
		return 0x02, i[1:4]
	}

	return 0x03, i
}

// Returns a byte array corresponding to the "Update Stream" format
func BuildUpdateStream(dataType byte, index uint32, payload []byte) []byte {
	var header byte = 0x00

	lenI, i := CompactIndex(index)

	header |= lenI
	header |= (dataType << 2)

	output := make([]byte, 0, 1+len(i)+len(payload))
	output[0] = header
	output = append(output, i...)
	output = append(output, payload...)

	return output
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
