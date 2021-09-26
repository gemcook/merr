package errs

import (
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
