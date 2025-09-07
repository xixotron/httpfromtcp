package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespaces
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 37, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	headers = Headers{"host": "localhost:42069"}
	data = []byte("User-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "curl/7.81.0", headers["user-agent"])
	assert.False(t, done)
	assert.Equal(t, 25, n)
	n, done, err = headers.Parse(data[n:])
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 13, n)
	assert.False(t, done)
	assert.Equal(t,
		Headers{
			"host":       "localhost:42069",
			"user-agent": "curl/7.81.0",
			"accept":     "*/*",
		},
		headers,
	)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("\r\nBody content")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Invalid single header name
	headers = NewHeaders()
	data = []byte("H@st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
