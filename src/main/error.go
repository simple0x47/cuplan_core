package main

type Error struct {
	errorKind string
	message   string
}

func NewError(errorKind string, message string) *Error {
	e := new(Error)
	e.errorKind = errorKind
	e.message = message

	return e
}

func (e Error) ErrorKind() string {
	return e.errorKind
}

func (e Error) Message() string {
	return e.message
}
