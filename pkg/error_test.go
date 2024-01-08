package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestError_ErrorKind(t *testing.T) {
	errorKind := "error_kind"
	newError := NewError(errorKind, "")

	assert.Equal(t, errorKind, newError.ErrorKind())
}

func TestError_Message(t *testing.T) {
	message := "message"
	newError := NewError("", message)

	assert.Equal(t, message, newError.message)
}
