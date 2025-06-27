package helpers

import (
	"bytes"
	"encoding/gob"
)

// Serialize encodes a value using gob.
func Serialize(src any) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(src); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Deserialize decodes a value using gob.
func Deserialize(src []byte, dst any) error {
	return gob.NewDecoder(bytes.NewBuffer(src)).Decode(dst)
}
