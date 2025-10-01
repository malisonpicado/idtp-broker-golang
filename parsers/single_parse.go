package parsers

import (
	"idtp/utils"
	"idtp/values"
)

type Request struct {
	Method  values.MethodCode
	Type    values.DataType
	Index   uint32
	Payload []byte
}

func ExtractIndexLength(header byte) byte {
	return (header & 0x03) + 1
}

func ExtractDeviceType(header byte) byte {
	return (header >> 2) & 0x0F
}

func ExtractMethod(header byte) byte {
	return header >> 7
}

// Parse an array of N bytes into a single request, following protocol
// request format. Returns a request struct and the number of
// bytes evaluated. If number of bytes evaluated is 0, means error;
// only EXPAND method returns the value 0xFF.
//
// Returns 0 only if payload is incomplete or unable to evaluate;
// this function does not validate data, just parse it.
func SingleParse(input []byte, useExtendedMode bool) (req Request, bytesEvaluated byte) {
	if len(input) == 0 {
		return Request{}, 0
	}

	start := byte(0)
	method := byte(0)

	if useExtendedMode {
		start = 1
		method = input[0] & 0x07
	} else {
		method = ExtractMethod(input[0])
	}

	switch values.MethodCode(method) {
	case values.GET:
		return ParseGet(input[start:])
	case values.UPDATE:
		return ParseUpdate(input[start:])
	case values.EXPAND:
		return ParseExpand(input[start:]), 0xFF
	case values.SET_TYPE:
		return ParseSetType(input[start:])
	}

	return Request{}, 0
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

func ParseGet(input []byte) (req Request, bytesEvaluated byte) {
	ilen := ExtractIndexLength(input[0])

	if len(input) < int(1+ilen) {
		return Request{}, 0
	}

	index := utils.BytesToU32(input[1 : 1+ilen])

	return Request{
		Method: values.GET,
		Type:   values.INT32,
		Index:  index,
	}, (1 + ilen)
}

func ParseUpdate(input []byte) (req Request, bytesEvaluated byte) {
	ilen := ExtractIndexLength(input[0])
	dtype := ExtractDeviceType(input[0])
	dsize := SizeOf(dtype)

	l := 1 + ilen + SizeOf(dtype)

	if len(input) < int(l) || dsize == 0 {
		return Request{}, 0
	}

	index := utils.BytesToU32(input[1 : 1+ilen])
	payload := input[1+ilen : l]

	return Request{
		Method:  values.UPDATE,
		Type:    values.DataType(dtype),
		Index:   index,
		Payload: payload,
	}, l
}

func ParseExpand(input []byte) Request {
	return Request{
		Method:  values.EXPAND,
		Type:    values.INT32,
		Index:   0,
		Payload: input,
	}
}

func ParseSetType(input []byte) (req Request, bytesEvaluated byte) {
	ilen := ExtractIndexLength(input[0])
	index := utils.BytesToU32(input[1 : 1+ilen])

	if len(input) < int(2+ilen) {
		return Request{}, 0
	}

	return Request{
		Method:  values.SET_TYPE,
		Type:    values.INT32,
		Index:   index,
		Payload: input[1+ilen : 2+ilen],
	}, (2 + ilen)
}
