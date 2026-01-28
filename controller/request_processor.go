package controller

import (
	"idtp/parsers"
	"idtp/storage"
	"idtp/utils"
	"idtp/values"
	"net"
	"strings"
)

// From an array of bytes that encodes a N block of requests, the function parses
// the request and executes the request. Returns a response in the same order of request
// encoded in bytes. The request param, must be cleaned before the use of this
// function; clean from the first byte of command that enters from a entity input buffer.
func RequestProcessor(
	request []byte,
	config values.Configuration,
	currentEntity net.Conn,
	entityConfig *values.Entity,
	storage *storage.Storage,
	clients *storage.ClientsList,
	dependents *storage.DependentsManager) []byte {

	var freeAllow bool = !(entityConfig != nil && entityConfig.ProcessAsStrict)
	var useExtended bool = entityConfig != nil && entityConfig.EntityType == values.ENTITY_CLIENT

	requests, _ := parsers.ParseRequest(request, useExtended)
	var responses []byte

	for _, req := range requests {
		if req.Method == values.GET {
			if !freeAllow && !utils.HasIndex(req.Index, entityConfig.DependencyParams) {
				responses = append(responses, values.RC_VARIABLE_OPERATION_NOT_ALLOWED)
				continue
			}

			responses = append(responses, storage.GetAt(req.Index)...)
			continue
		}

		if req.Method == values.UPDATE {
			if !freeAllow && !utils.HasIndex(req.Index, entityConfig.DependencyParams) {
				responses = append(responses, values.RC_VARIABLE_OPERATION_NOT_ALLOWED)
				continue
			}

			response := storage.UpdateAt(req.Index, byte(req.DataType), req.Payload, freeAllow)
			responses = append(responses, response)

			if response == values.RC_SUCCESS {
				payload := append([]byte{0xFF}, parsers.BuildUpdateStream(req.DataType, req.Index, req.Payload)...)
				go Broadcast(req.Index, payload, dependents, clients, currentEntity)
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
				payload := append([]byte{0xFF}, parsers.BuildUpdateStream(req.DataType, req.Index, make([]byte, utils.SizeOf(byte(req.DataType))))...)
				go Broadcast(req.Index, payload, dependents, clients, currentEntity)
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
func ConnectionRequestProcessor(request []byte, config values.Configuration, stg *storage.Storage) (connReq values.ConnectionRequest, code byte) {
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

	// Verify parameter is valid
	for _, param := range connreq.Parameters {
		if !stg.IndexExists(param.Index) {
			return values.ConnectionRequest{}, values.RC_INVALID_PARAMETER
		}
	}

	return connreq, values.RC_SUCCESS
}
