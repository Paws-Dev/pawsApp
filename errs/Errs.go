package errs

import (
	"errors"
	"strings"
)

type Error struct {
	error   error
	handler func(error) error
}

func New(message string) *Error {
	return &Error{error: errors.New(message)}

}

func FromError(err error) *Error {
	if err == nil {
		return &Error{}
	}
	return &Error{
		error: err,
	}
}

func (e *Error) Root() string {
	if e == nil || e.error == nil {
		return ""
	}
	trace := strings.Split(e.error.Error(), "\n")
	if len(trace) == 0 {
		return ""
	}
	return trace[len(trace)-1]
}

func (e *Error) ErrorF(symbol string) string {
	if e == nil || e.error == nil {
		return ""
	}
	replacer := strings.NewReplacer(
		"\r\n", symbol,
		"\n", symbol,
	)
	return replacer.Replace(e.error.Error())
}

func (e *Error) Error() string {
	return e.error.Error()
}

func (e *Error) Add(err error) *Error {
	if e == nil {
		return nil
	}

	if err == nil {
		return e
	}

	if e.error == nil {
		e.error = err
	} else {
		e.error = errors.Join(e.error, err)
	}

	return e
}

func (e *Error) Is(target error) bool {
	if e == nil || e.error == nil {
		return false
	}
	return errors.Is(e.error, target)
}

func (e *Error) As(target interface{}) bool {
	if e == nil || e.error == nil {
		return false
	}
	return errors.As(e.error, &target)
}
