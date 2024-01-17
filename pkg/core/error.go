package core

import "fmt"

type Error struct {
	ErrorKind string `json:"error_kind"`
	Message   string `json:"message"`
}

func NewError(errorKind string, message string) *Error {
	e := new(Error)
	e.ErrorKind = errorKind
	e.Message = message

	return e
}

func (e Error) String() string {
	return fmt.Sprintf("%s: %s", e.ErrorKind, e.Message)
}
