package httpserver

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestDecodeEncodedAvatar(t *testing.T) {
	code := GetEncodedAvatar("41253fa3-1761-4b53-8664-2e0daf010f20.png", 1)
	t.Log("code", code)
	fileName, userID, err := DecodeEncodedAvatar(code)
	assert.NilError(t, err)
	assert.Equal(t, userID, uint32(54))
	assert.Equal(t, fileName, "41253fa3-1761-4b53-8664-2e0daf010f20.png")
}
