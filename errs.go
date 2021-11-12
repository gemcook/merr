package merr

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"sync"
	"text/tabwriter"
)

// tabwriter parameter
const (
	minWidth  int    = 2
	tabwidth  int    = 0
	padding   int    = 1
	padchar   byte   = ' '
	delimiter string = ","
)

var (
	defaultOutput  io.Writer = os.Stdout
	output         io.Writer = defaultOutput
	defaultNewline string    = "\u000A"
	newLine        string    = defaultNewline
)

var formatter = func(e *errs) string {
	var result string
	suffix := newLine
	for i, err := range e.Errors {
		if len(e.Errors)-1 == i {
			suffix = ""
		}
		result += err.Error() + suffix
	}
	return result
}

func GetOutput() io.Writer {
	return output
}

func SetOutput(out io.Writer) {
	output = out
}

func ResetOutput() {
	output = defaultOutput
}

// SetNewLine is a function to set the newline code.
//
// Only "\n" , "\r\n" or "\r" are allowed.
// If they do not apply, the default value ("\n") is used.
func SetNewLine(n string) {
	if n == "\u000A" || n == "\u000D\u000A" || n == "\u000D" {
		newLine = n
		return
	}
	newLine = defaultNewline
}

type Errs interface {
	Append(err error)
	Error() string
	Is(target error) bool
	As(target interface{}) bool
	PrettyPrint()
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

func (e *errs) As(target interface{}) bool {
	for _, err := range e.Errors {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}

func (e *errs) PrettyPrint() {
	fmt.Fprint(output, e.prettyFormat())
}

type prettyPrinter struct {
	buf   *bytes.Buffer
	tw    *tabwriter.Writer
	depth int
}

func newPrettyPrinter(depth int) *prettyPrinter {
	p := &prettyPrinter{
		buf:   bytes.NewBufferString(""),
		depth: depth,
	}
	p.tw = tabwriter.NewWriter(p.buf, minWidth, tabwidth, padding, padchar, 0)
	return p
}

func (e *errs) prettyFormat() string {
	p := newPrettyPrinter(0)
	fmt.Fprint(p.tw, "Errors[")

	p.depth++
	for _, err := range e.Errors {
		fmt.Fprint(p.tw, newLine)
		p.writeValue(reflect.ValueOf(err), true)
		fmt.Fprint(p.tw, delimiter)
	}
	p.depth--

	if len(e.Errors) > 0 {
		fmt.Fprint(p.tw, newLine)
	}
	fmt.Fprint(p.tw, "]")
	p.tw.Flush()
	return p.buf.String()
}

func (p *prettyPrinter) indent() {
	indent := strings.Repeat("\t", p.depth)
	fmt.Fprint(p.tw, indent)
}

func (p *prettyPrinter) writeValue(val reflect.Value, enableIndent bool) {
	if enableIndent {
		p.indent()
	}

	switch val.Kind() {
	case reflect.Bool:
		fmt.Fprintf(p.tw, "%#v", val.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(p.tw, "%#v", val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(p.tw, "%#v", val.Uint())
	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(p.tw, "%#v", val.Float())
	case reflect.Complex64, reflect.Complex128:
		fmt.Fprintf(p.tw, "%#v", val.Complex())
	case reflect.String:
		fmt.Fprintf(p.tw, "%#v", val.String())
	case reflect.Map:
		fmt.Fprintf(p.tw, "%s{", val.Type().String())
		if !val.IsNil() {
			fmt.Fprint(p.tw, newLine)
			keys := val.MapKeys()
			p.depth++
			for i := range keys {
				p.writeValue(keys[i], true)
				fmt.Fprint(p.tw, ":")

				mapValuePrinter := newPrettyPrinter(p.depth)
				mapValuePrinter.writeValue(val.MapIndex(keys[i]), false)
				mapValuePrinter.tw.Flush()

				fmt.Fprint(p.tw, mapValuePrinter.buf.String())
				fmt.Fprint(p.tw, delimiter)
				fmt.Fprint(p.tw, newLine)
			}
			p.depth--
		}
		p.indent()
		fmt.Fprint(p.tw, "}")
	case reflect.Struct:
		fmt.Fprint(p.tw, val.Type().String())
		if val.NumField() == 0 {
			fmt.Fprint(p.tw, "{}")
			return
		}
		fmt.Fprint(p.tw, "{")
		if val.IsValid() {
			fmt.Fprint(p.tw, newLine)
			p.depth++
			for i := 0; i < val.NumField(); i++ {
				p.indent()
				fmt.Fprintf(p.tw, "%s:\t", val.Type().Field(i).Name)

				structValuePrinter := newPrettyPrinter(p.depth)
				structValuePrinter.writeValue(val.Field(i), false)
				structValuePrinter.tw.Flush()

				fmt.Fprint(p.tw, structValuePrinter.buf.String())
				fmt.Fprint(p.tw, delimiter)
				fmt.Fprint(p.tw, newLine)
			}
			p.depth--
		}
		p.indent()
		fmt.Fprint(p.tw, "}")
	case reflect.Interface:
		switch elm := val.Elem(); {
		case elm.Kind() == reflect.Invalid:
			fmt.Fprint(p.tw, "nil")
		case elm.IsValid():
			p.writeValue(elm, false)
		default:
			fmt.Fprint(p.tw, val.Type().String())
			fmt.Fprint(p.tw, "nil")
		}
	case reflect.Array, reflect.Slice:
		fmt.Fprint(p.tw, val.Type().String())
		if val.Kind() == reflect.Slice && val.IsNil() {
			fmt.Fprint(p.tw, "(nil)")
			return
		}
		if val.Len() == 0 {
			fmt.Fprint(p.tw, "{}")
			return
		}
		fmt.Fprint(p.tw, "{")
		if val.IsValid() {
			fmt.Fprint(p.tw, newLine)
			p.depth++
			for i := 0; i < val.Len(); i++ {
				p.indent()
				p.writeValue(val.Index(i), false)
				fmt.Fprint(p.tw, delimiter)
				fmt.Fprint(p.tw, newLine)
			}
			p.depth--
		}
		p.indent()
		fmt.Fprint(p.tw, "}")
	case reflect.Ptr:
		elm := val.Elem()
		if elm.IsValid() {
			fmt.Fprint(p.tw, "&")
			p.writeValue(elm, false)
		} else {
			fmt.Fprint(p.tw, "(&"+val.Type().Name()+")(nil)")
		}
	case reflect.Chan:
		fmt.Fprintf(p.tw, "%s(%#v)", val.Type().String(), val.Pointer())
	case reflect.Func:
		fmt.Fprint(p.tw, val.Type().String()+" {...}")
	case reflect.Invalid:
		fmt.Fprint(p.tw, "nil")
	}
}
