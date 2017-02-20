package lib

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionedJoinAndSplit(t *testing.T) {
	testBytesIn := []byte("abcdefghijklmnopqrstuvwxyz")
	testBytesVersion := make([]byte, 4)
	binary.LittleEndian.PutUint32(testBytesVersion, Version)
	testBytesOut := append(testBytesVersion, testBytesIn...)

	joined := VersionedJoin(testBytesIn)
	assert.True(t, bytes.Equal(joined, testBytesOut), "conversion should match")

	version, ciphertext, nonce, salt, err := VersionedSplit(joined)
	assert.Nil(t, err)
	assert.Equal(t, Version, version)
	assert.Equal(t, "abcdefghijkl", string(salt))
	assert.Equal(t, "mnopqrstuvwx", string(nonce))
	assert.Equal(t, "yz", string(ciphertext))
}
