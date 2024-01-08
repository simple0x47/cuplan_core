package secret

import (
	core "github.com/simpleg-eu/cuplan-core/pkg"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var provider *BitwardenProvider

func TestMain(m *testing.M) {
	setup()

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestBitwardenProvider_Get_ValidSecretId_ReturnsExpectedSecret(t *testing.T) {
	const secret = "le_secret :)"
	const secretId = "7c1d5dfd-a58b-47cf-bee5-b0a600fe50c9"

	result := provider.Get(secretId)

	assert.True(t, result.IsOk())
	assert.Equal(t, secret, result.Unwrap())
}

func TestBitwardenProvider_Get_InvalidSecretId_ReturnsError(t *testing.T) {
	const invalidSecretId = "1234"

	result := provider.Get(invalidSecretId)

	assert.False(t, result.IsOk())
	assert.Equal(t, core.CommandFailure, result.UnwrapErr().ErrorKind())
}

func setup() {
	provider = NewBitwardenProvider(os.Getenv("SECRETS_MANAGER_ACCESS_TOKEN"))
}
