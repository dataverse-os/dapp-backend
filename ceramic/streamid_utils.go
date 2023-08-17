package ceramic

import (
	"encoding/binary"
	"fmt"
)

// GetUVarInt returns the uint, rest of the buffer, bytes consumed, and error
// https://github.com/multiformats/unsigned-varint
func GetUVarInt(input []byte) (uint64, []byte, int, error) {
	code, idx := binary.Uvarint(input)
	if idx <= 0 {
		return 0, input, idx, fmt.Errorf("unable to unpack unsigned var int %v", input)
	}
	return code, input[idx:], idx, nil
}

func PutUVarInt(code uint64) []byte {
	buf := make([]byte, 10)
	n := binary.PutUvarint(buf, uint64(code))
	return buf[:n]
}
