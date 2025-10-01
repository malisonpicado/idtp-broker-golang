package storage

import (
	"idtp/parsers"
	"idtp/utils"
	"idtp/values"
	"sync"
)

type Slot struct {
	Type       values.DataType
	Payload    []byte
	Dependents map[uint32]struct{}
}

type Storage struct {
	Length uint32
	Limit  uint32
	Data   []Slot
	Mu     sync.RWMutex
}

// Builds a variable slot. If data dype is not valid, it
// set the default data type (INT32)
func SlotBuilder(datatype uint8) (Slot, values.StatusCode) {
	dataType := values.DataType(datatype)
	status := values.SUCCESS

	if dataType > values.FLOAT64 {
		dataType = values.INT32
		status = values.INVALID_DATA_TYPE_SET_TO_DEFAULT
	}

	return Slot{
		Type:       dataType,
		Payload:    make([]byte, values.SizeOf(dataType)),
		Dependents: make(map[uint32]struct{}),
	}, status
}

// Initializes and returns an empty storage
func Initialize(config values.Configuration) *Storage {
	return &Storage{
		Limit:  config.StorageLimit,
		Length: 0,
		Data:   make([]Slot, 0, config.StorageLimit),
		Mu:     sync.RWMutex{},
	}
}

// Adds new variables to the storage. The number of variables added is proportional to
// the size of the slots parameter. The slots parameter of the function is an array of
// data types; if the data type is invalid, the resulting slot will have the default
// data type. The function returns a byte array according to the protocol specification
// for processing the EXPAND method.
func (storage *Storage) Expand(slots []byte) []byte {
	storage.Mu.Lock()
	defer storage.Mu.Unlock()

	if len(slots) == 0 {
		return nil
	}

	currentLength := storage.Length
	limit := storage.Limit

	data := make([]Slot, 0, len(slots))
	result := make([]byte, 0, 5*len(slots))

	for i, item := range slots {
		index := currentLength + uint32(i)

		if index >= limit {
			result = append(result, byte(values.EXPANSION_LIMIT_REACHED))
			continue
		}

		slot, statusCode := SlotBuilder(item)
		result = append(result, byte(statusCode))
		result = append(result, utils.U32ToBytes(index)...)
		data = append(data, slot)
	}

	storage.Length += uint32(len(data))
	storage.Data = append(storage.Data, data...)
	return result
}

// Gets the value of a variable at index I. Returns an
// UpdateStream format in bytes with its status code at start.
// If error returns a single byte of operation status code.
func (storage *Storage) GetAt(index uint32) []byte {
	storage.Mu.RLock()
	defer storage.Mu.RUnlock()

	if index >= storage.Length {
		return []byte{byte(values.INVALID_INDEX)}
	}

	slot := storage.Data[index]
	upst := parsers.BuildUpdateStream(byte(slot.Type), index, slot.Payload)

	result := make([]byte, 0, 1+len(upst))
	result = append(result, byte(values.SUCCESS))
	result = append(result, upst...)

	return result
}

// Updates the value of a variable at index I. Returns a single byte slice
// that reprent the update operation status code.
func (storage *Storage) UpdateAt(index uint32, datatype byte, payload []byte, config values.Configuration) []byte {
	storage.Mu.Lock()
	defer storage.Mu.Unlock()

	if index >= storage.Length {
		return []byte{byte(values.INVALID_INDEX)}
	}

	if values.SizeOf(values.DataType(datatype)) == 0 {
		return []byte{byte(values.UNKNOWN_DATA_TYPE)}
	}

	t := storage.Data[index]

	if config.OperationMode != values.OP_MODE_STRICT {
		t.Type = values.DataType(datatype)
	} else if values.DataType(datatype) != t.Type {
		return []byte{byte(values.DATA_TYPE_OVERWRITE_NOT_ALLOWED)}
	}

	t.Payload = payload
	storage.Data[index] = t

	return []byte{byte(values.SUCCESS)}
}

// NOTA: Se entiende que este método es usado cuando el tipo de entidad
// es "cliente" y ha sido previamente validada
func (storage *Storage) SetTypeAt(index uint32, newType byte) byte {
	storage.Mu.Lock()
	defer storage.Mu.Unlock()

	if index >= storage.Length {
		return byte(values.INVALID_INDEX)
	}

	slot, statusCode := SlotBuilder(newType)

	t := storage.Data[index]

	t.Type = slot.Type
	t.Payload = slot.Payload

	storage.Data[index] = t

	return byte(statusCode)
}
