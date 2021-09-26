package errs

import (
	"errors"
	"sync"
)

var formatter = func(e *errs) string {
	var result string
	for _, e := range e.Errors {
		result += e.Error()
	}
	return result
}

type Errs interface {
	Append(err error)
	Error() string
	Is(target error) bool
}

type errs struct {
	mx     sync.Mutex
	Errors []error
}

func New() Errs {
	return &errs{
		mx:     sync.Mutex{},
		Errors: nil,
	}
}

func (e *errs) Error() string {
	return formatter(e)
}

func (e *errs) Append(err error) {
	e.mx.Lock()
	defer e.mx.Unlock()
	if e.Errors == nil {
		e.Errors = make([]error, 0)
	}
	e.Errors = append(e.Errors, err)
}

func (e *errs) Is(target error) bool {
	for _, err := range e.Errors {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}
