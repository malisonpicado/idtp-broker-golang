package parsers

import (
	"idtp/utils"
)

type UpdateStream struct {
	DataType uint8
	Index    uint32
	Payload  []byte
}

// Returns the size in bytes of a value that represents a
// data type. If the value t does not represents a valid
// data type, then the function will return 0
func SizeOf(t uint8) uint8 {
	switch t {
	case 0, 1, 5:
		return 1
	case 2, 6:
		return 2
	case 3, 7, 9:
		return 4
	case 4, 8, 10:
		return 8
	}

	return 0
}

// Converts an index to byte format. Returns the index length
// code and the index in bytes based on the length code value.
// 0x00 -> 1 byte, ..., 0x03 -> 4 bytes
func CompactIndex(index uint32) (ilen byte, ibyte []byte) {
	i := utils.U32ToBytes(index)

	if index < 256 {
		return 0x00, i[3:4]
	}

	if index < 65_535 {
		return 0x01, i[2:4]
	}

	if index < 16_777_215 {
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

	output = append(output, header)
	output = append(output, i...)
	output = append(output, payload...)

	return output
}

// Returns an UpdateStream struct from an array of N bytes.
// Returns 0 if payload is missing or unable to evaluate.
func ParseUpdateStream(data []byte) (upstr UpdateStream, bytesEvaluated byte) {
	ilen := ExtractIndexLength(data[0])
	dtype := ExtractDeviceType(data[0])
	dsize := SizeOf(dtype)
	l := 1 + ilen + dsize

	if len(data) < int(l) || dsize == 0 {
		return UpdateStream{}, 0
	}

	index := utils.BytesToU32(data[1 : 1+ilen])

	return UpdateStream{
		DataType: dtype,
		Index:    index,
		Payload:  data[1+ilen:],
	}, l
}
