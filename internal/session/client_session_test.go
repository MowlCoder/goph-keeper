package session

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientSession_SetToken(t *testing.T) {
	file, err := os.Create(t.TempDir() + "/test.json")
	defer file.Close()
	require.NoError(t, err)
	session := NewClientSession(file)
	token := "test-token"

	session.SetToken(token)
	assert.Equal(t, token, session.Token)
}

func TestClientSession_IsAuth(t *testing.T) {
	file, err := os.Create(t.TempDir() + "/test.json")
	defer file.Close()
	require.NoError(t, err)
	session := NewClientSession(file)
	token := "test-token"

	assert.Equal(t, false, session.IsAuth())
	session.SetToken(token)
	assert.Equal(t, token, session.Token)
	assert.Equal(t, true, session.IsAuth())
}
