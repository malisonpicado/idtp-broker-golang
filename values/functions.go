package values

// Returns the size in bytes of a data type
// payload. If the returned size is 0, the
// the data type is invalid.
func SizeOf(dataType byte) byte {
	switch dataType {
	case BOOLEAN, UINT8, INT8:
		return 1
	case UINT16, INT16:
		return 2
	case UINT32, INT32, FLOAT32:
		return 4
	case UINT64, INT64, FLOAT64:
		return 8
	}

	return 0
}
