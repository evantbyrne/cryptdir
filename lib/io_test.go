package lib

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptFileAndDecryptFile(t *testing.T) {
	var (
		cleartext1 []byte
		cleartext2 []byte
		key1       []byte
		key2       []byte
		nonce1     []byte
		nonce2     []byte
		password   []byte
		salt1      []byte
		salt2      []byte
	)

	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatal(err)
	}

	cleartext1 = []byte("The quick brown fox jumps over the lazy dog.")
	password = []byte("mysuperdupersecurepassword")
	salt1 = NewNonce()
	key1 = NewKey(salt1, password)
	nonce1 = NewNonce()

	EncryptFile(tmp.Name(), cleartext1, nonce1, key1, salt1)
	cleartext2, nonce2, key2, salt2 = DecryptFile(tmp.Name(), password)

	assert.Equal(t, cleartext1, cleartext2, "cleartext should match")
	assert.Equal(t, nonce1, nonce2, "nonce should match")
	assert.Equal(t, key1, key2, "key should match")
	assert.Equal(t, salt1, salt2, "salt should match")
}

func TestGetRawFileName(t *testing.T) {
	var (
		config     Config
		filenameA1 string
		filenameA2 string
		filenameB1 string
	)

	config.Files = make(map[string]string, 0)

	filenameA1 = getRawFileName(config, "hello.txt")
	assert.Len(t, filenameA1, rawFileNameLength, "generated raw filename is incorrect length")
	config.Files["hello.txt"] = filenameA1

	filenameA2 = getRawFileName(config, "hello.txt")
	assert.Equal(t, filenameA1, filenameA2, "raw filenames should match")

	filenameB1 = getRawFileName(config, "foo.jpg")
	assert.NotEqual(t, filenameA1, filenameB1, "raw filenames should not match")
}
