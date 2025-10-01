package parsers

import "idtp/values"

func BuildGetRequest(index uint32) []byte {
	header := byte(0x00)
	ilen, ci := CompactIndex(index)

	header |= ilen
	t := make([]byte, 2+ilen)

	t = append(t, header)
	t = append(t, ci...)

	return t
}

func BuildUpdateRequest(index uint32, datatype values.DataType, payload []byte) []byte {
	header := byte(0x80)
	dtype := datatype << 2
	ilen, ci := CompactIndex(index)

	header |= byte(dtype)
	header |= ilen

	reqLen := int(2+ilen) + len(ci) + len(payload)
	t := make([]byte, reqLen)

	t = append(t, header)
	t = append(t, ci...)
	t = append(t, payload...)

	return t
}

func BuildExpandedGetRequest(index uint32) []byte {
	single := BuildGetRequest(index)

	t := make([]byte, 1+len(single))

	t = append(t, 0x00)
	t = append(t, single...)

	return t
}

func BuildExpandedUpdateRequest(index uint32, datatype values.DataType, payload []byte) []byte {
	single := BuildUpdateRequest(index, datatype, payload)
	single[0] = single[0] & 0x7F

	t := make([]byte, 1+len(single))

	t = append(t, 0x01)
	t = append(t, single...)

	return t
}

func BuildExpandRequest(dtypes []byte) []byte {
	t := make([]byte, 1+len(dtypes))

	t = append(t, 0x02)
	t = append(t, dtypes...)

	return t
}

func BuildSetTypeRequest(index uint32, newdtype values.DataType) []byte {
	ilen, ci := CompactIndex(index)

	t := make([]byte, 3+ilen+1)

	t = append(t, 0x03, ilen)
	t = append(t, ci...)
	t = append(t, byte(newdtype))

	return t
}
