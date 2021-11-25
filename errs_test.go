package merr

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_errs_Error(t *testing.T) {
	tests := []struct {
		name string
		arg  error
		want string
	}{
		{
			"check/formatter",
			nil,
			"",
		},
		{
			"check/error",
			fmt.Errorf("%s", "error"),
			"error",
		},
		{
			"check/error/multiple",
			&errs{
				mx: sync.Mutex{},
				Errors: []error{
					fmt.Errorf("%s", "error1"),
					fmt.Errorf("%s", "error2"),
				},
			},
			"error1\nerror2",
		},
		{
			"check/error/wrapped",
			&errs{
				mx: sync.Mutex{},
				Errors: []error{
					fmt.Errorf("%w", &somethingError{}),
					fmt.Errorf("%w", &somethingError{}),
				},
			},
			"something error\nsomething error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			if tt.arg != nil {
				e.Append(tt.arg)
			}
			got := e.Error()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_errs_Is(t *testing.T) {
	fooErr := fmt.Errorf("%s", "foo")
	barErr := fmt.Errorf("%s", "bar")
	bazErr := fmt.Errorf("%s", "baz")

	tests := []struct {
		name string
		arg  error
		errs Errs
		want bool
	}{
		{
			"equal/contain",
			fooErr,
			&errs{
				mx: sync.Mutex{},
				Errors: []error{
					fooErr,
					barErr,
					bazErr,
				},
			},
			true,
		},
		{
			"equal/wrapped",
			barErr,
			&errs{
				mx: sync.Mutex{},
				Errors: []error{
					fooErr,
					fmt.Errorf("%w", barErr),
					bazErr,
				},
			},
			true,
		},
		{
			"not equal",
			bazErr,
			&errs{
				mx: sync.Mutex{},
				Errors: []error{
					fooErr,
					barErr,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.errs
			got := e.Is(tt.arg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_errs_As(t *testing.T) {
	serr := &somethingError{}

	tests := []struct {
		name string
		arg  interface{}
		errs Errs
		want bool
	}{
		{
			"match/contain",
			&serr,
			&errs{
				mx: sync.Mutex{},
				Errors: []error{
					&somethingError{},
					errors.New("err"),
				},
			},
			true,
		},
		{
			"match/wrapped",
			&serr,
			&errs{
				mx: sync.Mutex{},
				Errors: []error{
					fmt.Errorf("%w", &somethingError{}),
				},
			},
			true,
		},
		{
			"not match/not contain",
			&serr,
			&errs{
				mx: sync.Mutex{},
				Errors: []error{
					errors.New("err"),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.errs
			got := e.As(tt.arg)
			assert.Equal(t, tt.want, got)
		})
	}
}

type somethingError struct{}

func (s *somethingError) Error() string { return "something error" }

func Test_errs_PrettyPrint(t *testing.T) {
	tests := []struct {
		name string
		errs Errs
		out  string
	}{
		{
			"print struct",
			&errs{
				mx: sync.Mutex{},
				Errors: []error{
					&somethingError{},
				},
			},
			"Errors[\n  &merr.somethingError{},\n]",
		},
		{
			"print nil",
			&errs{
				mx: sync.Mutex{},
				Errors: []error{
					nil,
				},
			},
			"Errors[\n  nil,\n]",
		},
		{
			"print empty",
			New(),
			"Errors[]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.errs
			buf := bytes.NewBufferString("")
			SetOutput(buf)
			e.PrettyPrint()
			assert.Equal(t, tt.out, buf.String())
		})
	}
}

func TestSetNewLine(t *testing.T) {
	type args struct {
		n string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			`set \n`,
			args{n: "\n"},
			"Errors[\n  &errors.errorString{\n    s: \"foo\",\n  },\n]",
		},
		{
			`set \r\n`,
			args{n: "\r\n"},
			"Errors[\r\n  &errors.errorString{\r\n    s: \"foo\",\r\n  },\r\n]",
		},
		{
			`set \r`,
			args{n: "\r"},
			"Errors[\r &errors.errorString{\r   s: \"foo\",\r },\r]",
		},
		{
			`set \t but use default(\n)`,
			args{n: "\t"},
			"Errors[\n  &errors.errorString{\n    s: \"foo\",\n  },\n]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetNewLine(tt.args.n)

			e := New()
			err := fmt.Errorf("%s", "foo")
			e.Append(err)

			buf := bytes.NewBufferString("")
			SetOutput(buf)
			e.PrettyPrint()

			assert.Equal(t, tt.want, buf.String())
		})
	}
}
