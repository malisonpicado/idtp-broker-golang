package utils

// Returns the size in bytes of a data type
// payload. If the returned size is 0, the
// the data type is invalid.
func SizeOf(dataType byte) byte {
	switch dataType {
	// case of: boolean, uint8, int8
	case 0x00, 0x01, 0x05:
		return 1
	// case of: uint16, int16
	case 0x02, 0x06:
		return 2
	// case of: uint32, int32, float32
	case 0x03, 0x07, 0x09:
		return 4
	// case of: uint64, int64, float64
	case 0x04, 0x08, 0x0A:
		return 8
	}

	return 0
}

func HasIndex(i uint32, list []uint32) bool {
	for _, ix := range list {
		if ix == i {
			return true
		}
	}
	return false
}

// BytesToU16 converts a byte slice into an unsigned 16-bit integer (uint16).
// It takes a byte slice of any length and interprets the first two bytes as a uint16 value.
// If the input slice is shorter than 2 bytes, it pads the remaining bytes with zeros to form a complete uint16.
// The conversion uses little-endian format (least significant byte first).
func BytesToU16(arr []uint8) uint16 {
	var v uint16 = 0

	for i := range len(arr) {
		v <<= 8
		v |= uint16(arr[i])

		if i == 1 {
			break
		}
	}

	return v
}

// BytesToU32 converts a byte slice into an unsigned 32-bit integer (uint32).
// It takes a byte slice of any length and interprets the first four bytes as a uint32 value.
// If the input slice is shorter than 4 bytes, it pads the remaining bytes with zeros to form a complete uint32.
// The conversion uses little-endian format (least significant byte first).
func BytesToU32(arr []uint8) uint32 {
	var v uint32 = 0

	for i := range len(arr) {
		v <<= 8
		v |= uint32(arr[i])

		if i == 3 {
			break
		}
	}

	return v
}

// BytesToU64 converts a byte slice into an unsigned 64-bit integer (uint64).
// It takes a byte slice of any length and interprets the first eight bytes as a uint64 value.
// If the input slice is shorter than 8 bytes, it pads the remaining bytes with zeros to form a complete uint64.
// The conversion uses little-endian format (least significant byte first).
func BytesToU64(arr []uint8) uint64 {
	var v uint64 = 0

	for i := range len(arr) {
		v <<= 8
		v |= uint64(arr[i])

		if i == 7 {
			break
		}
	}

	return v
}

// U16ToBytes converts an unsigned 16-bit integer (uint16) into a byte slice.
// It uses little-endian format, meaning the least significant byte comes first.
// The resulting byte slice will always have a length of 2 bytes.
func U16ToBytes(v uint16) []byte {
	arr := make([]byte, 2)

	for i := range uint8(2) {
		arr[1-i] = uint8(v)
		v >>= 8
	}

	return arr
}

// U32ToBytes converts an unsigned 32-bit integer (uint32) into a byte slice.
// It uses little-endian format, meaning the least significant byte comes first.
// The resulting byte slice will always have a length of 4 bytes.
func U32ToBytes(v uint32) []byte {
	arr := make([]byte, 4)

	for i := range uint8(4) {
		arr[3-i] = uint8(v)
		v >>= 8
	}

	return arr
}

// U64ToBytes converts an unsigned 64-bit integer (uint64) into a byte slice.
// It uses little-endian format, meaning the least significant byte comes first.
// The resulting byte slice will always have a length of 8 bytes.
func U64ToBytes(v uint64) []byte {
	arr := make([]byte, 8)

	for i := range uint8(8) {
		arr[7-i] = uint8(v)
		v >>= 8
	}

	return arr
}
