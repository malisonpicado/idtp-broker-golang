package controller

import (
	"idtp/parsers"
	"idtp/storage"
	"idtp/values"
	"net"
	"strings"
)

// From an array of bytes that encodes a N block of requests, the function parses
// the request and executes the request. Returns a response in the same order of request
// encoded in bytes. The request param, must be cleaned before the use of this
// function; clean from the first byte of command that enters from a entity input buffer.
//
// useExtended bool comes from evaluating: "entity type" is "client"?
func RequestProcessor(
	request []byte,
	useExtended bool,
	config values.Configuration,
	currentEntity *net.Conn,
	entityConfig *values.Entity,
	storage *storage.Storage,
	clients *storage.ClientsList,
	dependents *storage.DependentsManager) []byte {

	requests, _ := parsers.ChainParse(request, useExtended)
	var responses []byte

	for _, req := range requests {
		if req.Method == values.GET {
			responses = append(responses, storage.GetAt(req.Index)...)
			continue
		}

		if req.Method == values.UPDATE {
			response := storage.UpdateAt(req.Index, byte(req.DataType), req.Payload, config.OperationMode == values.OP_MODE_STRICT, useExtended)
			responses = append(responses, response)

			if response == values.RC_SUCCESS {
				go Broadcast(req, []byte{0xFF, response}, dependents, clients, currentEntity)
			}

			continue
		}

		if req.Method == values.EXPAND {
			responses = append(responses, storage.Expand(req.Payload)...)
			continue
		}

		if req.Method == values.SET_TYPE {
			response := storage.SetTypeAt(req.Index, req.Payload[0])
			responses = append(responses, response)

			if response == values.RC_SUCCESS {
				payload := []byte{0xFF}
				payload = append(payload, parsers.BuildUpdateStream(byte(req.DataType), req.Index, make([]byte, values.SizeOf(byte(req.DataType))))...)
				go Broadcast(req, payload, dependents, clients, currentEntity)
			}

			continue
		}

		responses = append(responses, byte(values.RC_UNKNOWN_METHOD))
	}

	return responses
}

// From a array of bytes, parse the array information into a Connection Request.
// The request is validated, if not error, then the information provided
// is correct.
func ConnectionRequestProcessor(request []byte, config values.Configuration) (connReq values.ConnectionRequest, code byte) {
	connreq, statusCode := parsers.ConnectionRequestParse(request)

	if statusCode != values.RC_SUCCESS {
		return values.ConnectionRequest{}, statusCode
	}

	if connreq.ProtocolVersion != config.ProtocolVersion {
		return values.ConnectionRequest{}, values.RC_UNSUPPORTED_PROTOCOL_VERSION
	}

	if connreq.EntityType > values.ENTITY_CLIENT {
		return values.ConnectionRequest{}, values.RC_UNKNOWN_ENTITY_TYPE
	}

	if connreq.KeepAlive < 60 || connreq.KeepAlive > 3600 {
		return values.ConnectionRequest{}, values.RC_INVALID_KEEP_ALIVE
	}

	if config.OperationMode == values.OP_MODE_STRICT && len(connreq.Parameters) == 0 {
		return values.ConnectionRequest{}, values.RC_DEVICE_MUST_DECLARE_PARAMETERS
	}

	// If there is not authentication key, then any persistent
	// connection is allowed
	if len(config.Key) == 0 || config.OperationMode == values.OP_MODE_FREE {
		return connreq, values.RC_SUCCESS
	}

	// Default and strict mode requires any connection to authenticate
	if !connreq.HasAuth {
		return values.ConnectionRequest{}, values.RC_CONNECTION_MUST_AUTHENTICATE
	}

	// Mode is default or strict, and conn has auth
	if strings.Compare(connreq.UserKey, config.Key) != 0 {
		return values.ConnectionRequest{}, values.RC_FAILED_AUTHENTICATION
	}

	return connreq, values.RC_SUCCESS
}
