package storage

import (
	"idtp/parsers"
	"idtp/utils"
	"idtp/values"
	"sync"
)

type Variable struct {
	DataType byte
	Payload  []byte
	Mu       *sync.Mutex
}

type Storage struct {
	Data []Variable
	Mu   sync.Mutex
}

func VariableBuilder(datatype byte) (variable Variable, code byte) {
	code = values.RC_SUCCESS

	if datatype > values.FLOAT64 {
		datatype = values.INT32
		code = values.RC_INVALID_DATA_TYPE_SET_TO_DEFAULT
	}

	variable = Variable{
		DataType: datatype,
		Payload:  make([]byte, values.SizeOf(datatype)),
		Mu:       &sync.Mutex{},
	}

	return
}

func InitializeStorage(limit uint32) *Storage {
	return &Storage{
		Data: make([]Variable, 0, limit),
	}
}

// Expands the storage up to 'len(variables)'. Each
// byte is interpreted as the data type of the new
// variable, if wrong data type it will be set as the
// default data type.
//
// Returns the operation result and the assigned index
// of the new variable as a bytes for each new variable, as
// specified in IDTP standard.
//
//   - if success: [status code, index 0, index 1, index 2, index 3]
//   - if error: [status code]
func (storage *Storage) Expand(variables []byte) []byte {
	if len(variables) == 0 {
		return nil
	}

	storage.Mu.Lock()
	defer storage.Mu.Unlock()

	data := make([]Variable, 0, len(variables))
	result := make([]byte, 0, len(variables)*5)

	for i, item := range variables {
		index := len(storage.Data) + i

		if index >= cap(storage.Data) {
			result = append(result, values.RC_EXPANSION_LIMIT_REACHED)
			continue
		}

		slot, code := VariableBuilder(item)
		result = append(result, code)
		result = append(result, utils.U32ToBytes(uint32(index))...)
		data = append(data, slot)
	}

	storage.Data = append(storage.Data, data...)
	return result
}

// Gets the data of a variable at index "index".
//
// Returns the status code and the data in update stream format
// as specified in IDTP standard.
func (storage *Storage) GetAt(index uint32) []byte {
	if index >= uint32(len(storage.Data)) {
		return []byte{values.RC_INVALID_INDEX}
	}

	variable := &storage.Data[index]
	variable.Mu.Lock()
	defer variable.Mu.Unlock()

	data := make([]byte, 0)
	updateStream := parsers.BuildUpdateStream(variable.DataType, index, variable.Payload)
	data = append(data, values.RC_SUCCESS)
	data = append(data, updateStream...)

	return data
}

// Updates value and data type of a variable at index "index".
// It updates the data type if is not in strict mode or if entity is client.
//
// Returns the status code of the operation, as specified in IDTP standard.
func (storage *Storage) UpdateAt(index uint32, datatype byte, payload []byte, isStrictMode bool, isClient bool) byte {
	if index >= uint32(len(storage.Data)) {
		return values.RC_INVALID_INDEX
	}

	if datatype > values.FLOAT64 {
		return values.RC_UNKNOWN_DATA_TYPE
	}

	variable := &storage.Data[index]
	variable.Mu.Lock()
	defer variable.Mu.Unlock()

	isDifferentDataType := variable.DataType != datatype

	if isDifferentDataType && (isStrictMode && !isClient) {
		return values.RC_DATA_TYPE_OVERWRITE_NOT_ALLOWED
	}

	if isDifferentDataType {
		variable.DataType = datatype
	}

	variable.Payload = payload
	return values.RC_SUCCESS
}

// Sets a new data type of a variable at index "index" and updates its value
// to all zero bits, with a size according to the new data type.
//
// Returns the status code of the operation, as specified in IDTP standard.
func (storage *Storage) SetTypeAt(index uint32, newType byte) byte {
	if index >= uint32(len(storage.Data)) {
		return values.RC_INVALID_INDEX
	}

	if newType > values.FLOAT64 {
		return values.RC_UNKNOWN_DATA_TYPE
	}

	variable := &storage.Data[index]
	variable.Mu.Lock()
	defer variable.Mu.Unlock()

	variable.DataType = newType
	variable.Payload = make([]byte, values.SizeOf(newType))

	return values.RC_SUCCESS
}
