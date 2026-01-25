package parsers

import (
	"idtp/utils"
	"idtp/values"
)

func ParseRequest(buffer []byte, useExtendedMode bool) (req []values.Request, code byte) {
	next := 0
	var list []values.Request

	for next < len(buffer) {
		request, blockLen, e := parseSingleRequest(buffer[next:], useExtendedMode)

		// Invalid method or missing payload
		if e != values.RC_SUCCESS {
			return list, e
		}

		next += int(blockLen)
		list = append(list, request)

		// EXPAND method
		if blockLen == 0xFF {
			break
		}
	}

	return list, values.RC_SUCCESS
}

func parseSingleRequest(input []byte, useExtendedMode bool) (req values.Request, bytesEvaluated byte, error byte) {
	if len(input) == 0 {
		return values.Request{}, 0, values.RC_MISSING_PAYLOAD
	}

	var start byte = 0
	var method byte = 0

	if useExtendedMode {
		start = 1
		method = input[0] & 0x07
	} else {
		method = extractMethod(input[0])
	}

	switch method {
	case values.GET:
		return parseGet(input[start:])
	case values.UPDATE:
		return parseUpdate(input[start:])
	case values.EXPAND:
		return parseExpand(input[start:]), 0xFF, values.RC_SUCCESS
	case values.SET_TYPE:
		return parseSetType(input[start:])
	}

	return values.Request{}, 0, values.RC_UNKNOWN_METHOD
}

// =======================================
// NOTE on next functions:
// The functions always asume that the first byte is
// the request header, and header contains the data type
// and index length; except for ParseExpand, that takes all
// bytes as payload. So, the function SingleParse is
// responsible of parsing the request method and pass
// the correct input argument.
// =======================================

func parseGet(input []byte) (req values.Request, bytesEvaluated byte, error byte) {
	indexLength := extractIndexLength(input[0])
	blockSize := 1 + indexLength

	if len(input) < int(blockSize) {
		return values.Request{}, 0, values.RC_MISSING_PAYLOAD
	}

	index := utils.BytesToU32(input[1:blockSize])

	return values.Request{
		Method: values.GET,
		Index:  index,
	}, blockSize, values.RC_SUCCESS
}

func parseUpdate(input []byte) (req values.Request, bytesEvaluated byte, error byte) {
	indexLength := extractIndexLength(input[0])
	dataType := extractDataType(input[0])
	dataTypeSize := utils.SizeOf(dataType)
	blockSize := 1 + indexLength + dataTypeSize

	if len(input) < int(blockSize) {
		return values.Request{}, 0, values.RC_MISSING_PAYLOAD
	}

	if dataTypeSize == 0 {
		return values.Request{}, 0, values.RC_UNKNOWN_DATA_TYPE
	}

	return values.Request{
		Method:   values.UPDATE,
		DataType: dataType,
		Index:    utils.BytesToU32(input[1 : 1+indexLength]),
		Payload:  input[1+indexLength : blockSize],
	}, blockSize, values.RC_SUCCESS
}

func parseExpand(input []byte) values.Request {
	return values.Request{
		Method:  values.EXPAND,
		Payload: input,
	}
}

func parseSetType(input []byte) (req values.Request, bytesEvaluated byte, error byte) {
	indexLength := extractIndexLength(input[0])
	index := utils.BytesToU32(input[1 : 1+indexLength])
	blockSize := 2 + indexLength

	if len(input) < int(blockSize) {
		return values.Request{}, 0, values.RC_MISSING_PAYLOAD
	}

	return values.Request{
		Method:  values.SET_TYPE,
		Index:   index,
		Payload: input[1+indexLength : blockSize],
	}, blockSize, values.RC_SUCCESS
}

func extractIndexLength(header byte) byte {
	return (header & 0x03) + 1
}

func extractDataType(header byte) byte {
	return (header >> 2) & 0x0F
}

func extractMethod(header byte) byte {
	return header >> 7
}
