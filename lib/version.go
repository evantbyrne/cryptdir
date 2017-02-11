package lib

import (
	"encoding/binary"
	"errors"
)

const Version uint32 = 1

func VersionedJoin(in ...[]byte) []byte {
	out := make([]byte, 4)
	binary.LittleEndian.PutUint32(out, Version)
	for _, byteSlice := range in {
		out = append(out, byteSlice...)
	}

	return out
}

func VersionedSplit(in []byte) (version uint32, ciphertext, nonce, salt []byte, err error) {
	if len(in) < 25 {
		return 0, nil, nil, nil, errors.New("Invalid byte length.")
	}

	version = binary.LittleEndian.Uint32(in[:4])
	salt = in[4:16]
	nonce = in[16:28]
	ciphertext = in[28:]

	return version, salt, nonce, ciphertext, nil
}
