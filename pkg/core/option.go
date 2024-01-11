package core

// None is a constant representing there is no value for an Option value.
// DO NOT WRITE TO IT.
var None = Option[any]{
	value:    nil,
	hasValue: false,
}

// Option represents a value which may or may not have a value.
type Option[T any] struct {
	value    T
	hasValue bool
}

// IsSome checks whether the specified Option has a value.
func (o Option[any]) IsSome() bool {
	return o.hasValue
}

// IsNone checks whether the specified Option doesn't have a value.
func (o Option[any]) IsNone() bool {
	return !o.hasValue
}

// Some creates an Option which has a value.
func Some[T any](value T) Option[T] {
	return Option[T]{
		value:    value,
		hasValue: true,
	}
}

// Unwrap returns value if the Option has none, otherwise it panics.
func (o Option[T]) Unwrap() T {
	if o.IsNone() {
		panic("Tried to unwrap a 'None' value.")
	}

	return o.value
}
