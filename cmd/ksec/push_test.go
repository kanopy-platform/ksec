package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanSecretData(t *testing.T) {
	t.Parallel()

	reader := strings.NewReader(`
key=value
key1=
key2=value2
key3
`)

	data, err := readSecretData(reader)
	assert.NoError(t, err)

	assert.Equal(t, "value", string(data["key"]))
	assert.Equal(t, "", string(data["key1"]))
	assert.Equal(t, "value2", string(data["key2"]))

	_, ok := data["key3"]
	assert.False(t, ok)
}
