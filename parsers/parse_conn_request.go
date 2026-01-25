package parsers

import (
	"idtp/utils"
	v "idtp/values"
)

func ConnectionRequestParse(buffer []byte) (cr v.ConnectionRequest, code byte) {
	if len(buffer) < 5 {
		return v.ConnectionRequest{}, v.RC_MISSING_PAYLOAD
	}

	protocolVersion := buffer[0]
	entityType := buffer[1] & 0x07
	keepAlive := utils.BytesToU16(buffer[2:4])
	hasAuth := false
	authLen := buffer[4]
	key := ""

	if authLen > 0 {
		hasAuth = true

		if len(buffer) < int(5+authLen) {
			return v.ConnectionRequest{}, v.RC_MISSING_PAYLOAD
		}

		key = string(buffer[5 : 5+authLen])
	}

	var params []v.DeviceParameter
	next := int(5 + authLen)

	for next < len(buffer) {
		param, blockLen := ParseDeviceParameter(buffer[next:])

		if blockLen == 0 {
			return v.ConnectionRequest{}, v.RC_MISSING_PAYLOAD
		}

		next += int(blockLen)
		params = append(params, param)
	}

	return v.ConnectionRequest{
		ProtocolVersion: protocolVersion,
		EntityType:      entityType,
		KeepAlive:       keepAlive,
		HasAuth:         hasAuth,
		UserKey:         key,
		Parameters:      params,
	}, v.RC_SUCCESS
}

// Returns an EntityParameter from an array of N bytes.
// Return 0 for bytesEvaluated if payload is missing.
func ParseDeviceParameter(input []byte) (param v.DeviceParameter, bytesEvaluated byte) {
	paramType := extractMethod(input[0])
	indexLength := extractIndexLength(input[0])

	if len(input) < int(1+indexLength) {
		return v.DeviceParameter{}, 0
	}

	return v.DeviceParameter{
		Method: paramType,
		Index:  utils.BytesToU32(input[1 : 1+indexLength]),
	}, (1 + indexLength)
}
