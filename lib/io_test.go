package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRawFileName(t *testing.T) {
	var (
		config Config
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
