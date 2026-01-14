package parsers

import (
	"idtp/utils"
	"idtp/values"
)

func ConnectionRequestParse(buffer []byte) (cr values.ConnectionRequest, code byte) {
	if len(buffer) < 5 {
		return values.ConnectionRequest{}, values.RC_MISSING_PAYLOAD
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
			return values.ConnectionRequest{}, values.RC_MISSING_PAYLOAD
		}

		key = string(buffer[5 : 5+authLen])
	}

	var params []values.DeviceParameter
	next := int(5 + authLen)

	for next < len(buffer) {
		param, blockLen := ParseDeviceParameter(buffer[next:])

		if blockLen == 0 {
			return values.ConnectionRequest{}, values.RC_MISSING_PAYLOAD
		}

		next += int(blockLen)
		params = append(params, param)
	}

	return values.ConnectionRequest{
		ProtocolVersion: protocolVersion,
		EntityType:      entityType,
		KeepAlive:       keepAlive,
		HasAuth:         hasAuth,
		UserKey:         key,
		Parameters:      params,
	}, values.RC_SUCCESS
}

// Returns an EntityParameter from an array of N bytes.
// Return 0 for bytesEvaluated if payload is missing.
func ParseDeviceParameter(input []byte) (param values.DeviceParameter, bytesEvaluated byte) {
	paramType := extractMethod(input[0])
	indexLength := extractIndexLength(input[0])

	if len(input) < int(1+indexLength) {
		return values.DeviceParameter{}, 0
	}

	return values.DeviceParameter{
		Method: paramType,
		Index:  utils.BytesToU32(input[1 : 1+indexLength]),
	}, (1 + indexLength)
}

func EntityConfigBuilder(connreq values.ConnectionRequest) *values.Entity {
	entity := values.Entity{
		EntityType: connreq.EntityType,
	}

	for _, param := range connreq.Parameters {
		if param.Method == 0 {
			entity.DependencyParams = append(entity.DependencyParams, param.Index)
			continue
		}

		entity.UpdateParams = append(entity.UpdateParams, param.Index)
	}

	return &entity
}
