package parsers

func ChainParse(buffer []byte, useExtendedMode bool) []Request {
	next := uint32(0)
	list := []Request{}

	for next < uint32(len(buffer)) {
		request, blockLen := SingleParse(buffer[next:], useExtendedMode)

		// Invalid method or missing payload
		if blockLen == 0 {
			break
		}

		next += uint32(blockLen)
		list = append(list, request)

		// Expand method
		if blockLen == 0xFF {
			break
		}
	}

	return list
}
