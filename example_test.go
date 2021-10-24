package errs_test

import (
	"fmt"
	"os"

	"github.com/gemcook/errs"
)

type structError struct {
	i   int
	str string
	b   bool
}

func (structError) Error() string {
	return "structError"
}

type ptrError struct {
	i   int
	str string
	b   bool
}

func (*ptrError) Error() string {
	return "ptrError"
}

func Example() {
	err := errs.New()

	// error interface
	err.Append(fmt.Errorf("%w", fmt.Errorf("wrap error")))

	// struct
	var structErr structError = structError{
		i:   1,
		str: "error",
		b:   true,
	}
	err.Append(structErr)

	// ptr
	var ptrErr *ptrError = &ptrError{
		i:   1,
		str: "error",
		b:   true,
	}
	err.Append(ptrErr)

	errs.SetOutput(os.Stderr)
	err.PrettyPrint()
	errs.SetOutput(os.Stdout)
	err.PrettyPrint()

	// output:
	// Errors[
	//   &fmt.wrapError{
	//     msg: "wrap error",
	//     err: &errors.errorString{
	//       s: "wrap error",
	//     },
	//   },
	//   errs_test.structError{
	//     i:   1,
	//     str: "error",
	//     b:   true,
	//   },
	//   &errs_test.ptrError{
	//     i:   1,
	//     str: "error",
	//     b:   true,
	//   },
	// ]
}
