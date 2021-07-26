package errs

import (
	"fmt"
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
