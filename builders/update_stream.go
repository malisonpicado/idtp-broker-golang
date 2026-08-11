package builders

import "idtp/utils"

// Returns a byte array corresponding to the "Update Stream" format
func BuildUpdateStream(dataType byte, index uint32, payload []byte) []byte {
	var header byte = 0x00

	lenI, i := utils.CompactIndex(index)

	header |= lenI
	header |= (dataType << 2)

	output := make([]byte, 0, 1+len(i)+len(payload))
	output[0] = header
	output = append(output, i...)
	output = append(output, payload...)

	return output
}
