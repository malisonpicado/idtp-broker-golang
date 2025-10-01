package parsers

import (
	"idtp/utils"
	"idtp/values"
)

type EntityParameter struct {
	ParamType uint8
	Index     uint32
}

type ConnectionRequest struct {
	ProtocolVersion uint8
	EntityType      values.EntityType
	KeepAlive       uint16
	HasAuth         bool
	UserKey         string
	ParamUpdates    []uint32
	ParamDepends    []uint32
}

func ConnectionRequestParse(buffer []byte) (ConnectionRequest, values.StatusCode) {
	if len(buffer) < 5 {
		return ConnectionRequest{}, values.MISSING_PAYLOAD
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
			return ConnectionRequest{}, values.MISSING_PAYLOAD
		}

		key = string(buffer[5 : 5+authLen])
	}

	var upt []uint32
	var dep []uint32

	next := uint32(5 + authLen)

	for next < uint32(len(buffer)) {
		param, blockLen := ParseEntityParameter(buffer[next:])

		if blockLen == 0 {
			return ConnectionRequest{}, values.MISSING_PAYLOAD
		}

		next += uint32(blockLen)

		if param.ParamType == 0 {
			dep = append(dep, param.Index)
			continue
		}

		upt = append(upt, param.Index)
	}

	return ConnectionRequest{
		ProtocolVersion: protocolVersion,
		EntityType:      values.EntityType(entityType),
		KeepAlive:       keepAlive,
		HasAuth:         hasAuth,
		UserKey:         key,
		ParamUpdates:    upt,
		ParamDepends:    dep,
	}, values.SUCCESS
}

// Returns an EntityParameter from an array of N bytes.
// Return 0 for bytesEvaluated if payload is missing.
func ParseEntityParameter(input []byte) (param EntityParameter, bytesEvaluated byte) {
	paramType := ExtractMethod(input[0])
	ilen := ExtractIndexLength(input[0])

	if len(input) < int(1+ilen) {
		return EntityParameter{}, 0
	}

	return EntityParameter{
		ParamType: paramType,
		Index:     utils.BytesToU32(input[1 : 1+ilen]),
	}, (1 + ilen)
}
