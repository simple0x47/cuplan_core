package cuplancore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResult_Unwrap_OkValue(t *testing.T) {
	const expectedValue = "yes"
	result := Ok[string, int](expectedValue)

	value := result.Unwrap()

	assert.Equal(t, expectedValue, value)
}

func TestResult_Unwrap_ErrorValue_Panics(t *testing.T) {
	result := Err[string, int](1)

	assert.Panics(t, func() {
		result.Unwrap()
	})
}

func TestResult_UnwrapErr_OkValue_Panics(t *testing.T) {
	result := Ok[string, int]("abcd")

	assert.Panics(t, func() {
		result.UnwrapErr()
	})
}

func TestResult_UnwrapErr_ErrorValue(t *testing.T) {
	const expectedValue = "no"
	result := Err[int, string](expectedValue)

	value := result.UnwrapErr()

	assert.Equal(t, expectedValue, value)
}

func TestResult_IsOk_TrueIfOk(t *testing.T) {
	result := Ok[bool, bool](true)

	assert.True(t, result.IsOk())
}

func TestResult_IsOk_FalseIfError(t *testing.T) {
	result := Err[bool, bool](true)

	assert.False(t, result.IsOk())
}
