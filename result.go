package cuplancore

type Result[OkType any, ErrorType any] struct {
	ok   OkType
	err  ErrorType
	isOk bool
}

func Ok[OkType any, ErrorType any](ok OkType) Result[OkType, ErrorType] {
	return Result[OkType, ErrorType]{ok: ok, isOk: true}
}

func Err[OkType any, ErrorType any](err ErrorType) Result[OkType, ErrorType] {
	return Result[OkType, ErrorType]{err: err, isOk: false}
}

func (r Result[OkType, ErrorType]) IsOk() bool {
	return r.isOk
}

func (r Result[OkType, ErrorType]) Unwrap() OkType {
	if !r.isOk {
		panic("Unwrapped an 'Ok' when result contained an 'Err'.")
	}

	return r.ok
}

func (r Result[OkType, ErrorType]) UnwrapErr() ErrorType {
	if r.isOk {
		panic("Unwrapped an 'Err' when result contained an 'Ok'.")
	}

	return r.err
}
