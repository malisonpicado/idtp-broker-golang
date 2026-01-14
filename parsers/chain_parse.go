package parsers

import "idtp/values"

func ChainParse(buffer []byte, useExtendedMode bool) (req []values.Request, error byte) {
	next := 0
	var list []values.Request

	for next < len(buffer) {
		request, blockLen, e := SingleRequestParse(buffer[next:], useExtendedMode)

		// Invalid method or missing payload
		if e != values.RC_SUCCESS {
			return list, e
		}

		next += int(blockLen)
		list = append(list, request)

		// Expand method
		if blockLen == 0xFF {
			break
		}
	}

	return list, values.RC_SUCCESS
}
