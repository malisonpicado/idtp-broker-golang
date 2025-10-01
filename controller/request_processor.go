package controller

import (
	"idtp/parsers"
	"idtp/storage"
	"idtp/values"
	"strings"
)

// From an array of bytes that encodes a N block of requests, the function parses
// the request and executes the request. Returns a response in the same order of request
// encoded in bytes. The request param, must be cleaned before the use of this
// function; clean from the first byte of command that enters from a entity input buffer.
//
// useExtended bool comes from evaluating: "entity type" is "client"?
func RequestProcessor(request []byte, useExtended bool, storage *storage.Storage, config values.Configuration,
	entities *storage.EntitiesList,
	clients *storage.ClientsList, entityId uint32, omitSender bool) []byte {
	requests := parsers.ChainParse(request, useExtended)
	var updated []parsers.Request
	var responses []byte

	for _, req := range requests {
		if req.Method == values.GET {
			responses = append(responses, storage.GetAt(req.Index)...)
			continue
		}

		if req.Method == values.UPDATE {
			upt := storage.UpdateAt(req.Index, byte(req.Type), req.Payload, config)
			responses = append(responses, upt...)

			if upt[0] == byte(values.SUCCESS) {
				updated = append(updated, req)
			}

			continue
		}

		if req.Method == values.EXPAND {
			responses = append(responses, storage.Expand(req.Payload)...)
			continue
		}

		if req.Method == values.SET_TYPE {
			responses = append(responses, storage.SetTypeAt(req.Index, req.Payload[0]))
			continue
		}

		responses = append(responses, byte(values.UNKNOWN_METHOD))
	}

	go Broadcast(updated, entityId, omitSender, entities, clients, storage)
	return responses
}

// From a array of bytes, parse the array information into a Connection Request.
// The request is validated, if error is nil, then the information provided
// is correct.
func ConnectionRequestProcessor(request []byte, config values.Configuration) (parsers.ConnectionRequest, values.StatusCode) {
	connreq, statusCode := parsers.ConnectionRequestParse(request)

	if statusCode != values.SUCCESS {
		return parsers.ConnectionRequest{}, statusCode
	}

	if connreq.ProtocolVersion != config.ProtocolVersion {
		return parsers.ConnectionRequest{}, values.UNSUPPORTED_PROTOCOL_VERSION
	}

	if connreq.EntityType > values.ENTITY_CLIENT {
		return parsers.ConnectionRequest{}, values.UNKNOWN_ENTITY_TYPE
	}

	if connreq.KeepAlive < 60 || connreq.KeepAlive > 3600 {
		return parsers.ConnectionRequest{}, values.INVALID_KEEP_ALIVE
	}

	// If there is not authentication key, then any persistent
	// connection is allowed
	if len(config.Key) == 0 || config.OperationMode == values.OP_MODE_FREE {
		return connreq, values.SUCCESS
	}

	// Default and strict mode requires any connection to authenticate
	if !connreq.HasAuth && config.OperationMode != values.OP_MODE_FREE {
		return parsers.ConnectionRequest{}, values.CONNECTION_MUST_AUTHENTICATE
	}

	if connreq.HasAuth {
		if strings.Compare(connreq.UserKey, config.Key) != 0 {
			return parsers.ConnectionRequest{}, values.FAILED_AUTHENTICATION
		}
	}

	return connreq, values.SUCCESS
}
