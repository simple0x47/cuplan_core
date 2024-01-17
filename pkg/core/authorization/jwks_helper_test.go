package authorization

import (
	"github.com/simpleg-eu/cuplan-core/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type JwksHelperTestSuite struct {
	suite.Suite
}

func TestJwksHelperTestSuite(t *testing.T) {
	suite.Run(t, new(JwksHelperTestSuite))
}

func TestGetJwks_InvalidUrl_Error(t *testing.T) {
	const invalidUrl = "https://lmao"

	result := GetJwks(invalidUrl)

	assert.True(t, result.IsErr())
	assert.Equal(t, core.IOFailure, result.UnwrapErr().ErrorKind)
}

func TestGetJwks_NotAJwksUrl_Error(t *testing.T) {
	const notAJwksUrl = "https://json.org/example.html"

	result := GetJwks(notAJwksUrl)

	assert.True(t, result.IsErr())
	assert.Equal(t, core.SerializationFailure, result.UnwrapErr().ErrorKind)
}
