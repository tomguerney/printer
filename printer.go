package printer

import (
	"io"
	"os"

	"github.com/tomguerney/printer/internal/setter"
)

// Printer prints
type Printer struct {
	*setter.Setter
}

// NewSetter a new setter
func NewSetter() *setter.Setter {
	return setter.New()
}

var s = setter.New()

// OutWriter is the io.Writer to print to
var OutWriter io.Writer = os.Stdout

// ErrWriter is the io.Writer to print errors to
var ErrWriter io.Writer = os.Stderr

// SetOutWriter sets the OutWriter
func SetOutWriter(writer io.Writer) {
	s.OutWriter = writer
}

// SetErrWriter sets the ErrWriter
func SetErrWriter(writer io.Writer) {
	s.ErrWriter = writer
}

// Out prints formatted text to the OutWriter
func Out(i interface{}, a ...interface{}) {
	s.Out(i, a...)
}

// Err prints formatted text to ErrWriter
func Err(i interface{}, a ...interface{}) {
	s.Err(i, a...)
}

// Feed prints an empty line
func Feed() {
	s.Feed()
}

// Tabulate takes an array of string arrays and prints a table to output
func Tabulate(rows [][]string) {
	s.Tabulate(rows)
}

// TmplStencil applies a string map to the stencil with the passed ID and prints
// it to output
func TmplStencil(id string, data map[string]string) error {
	return s.TmplStencil(id, data)
}

// TableStencil take an array of string maps and prints stencilled rows to
// output
func TableStencil(id string, rows []map[string]string) error {
	return s.TableStencil(id, rows)
}

// AddTmplStencil adds a new template stencil
func AddTmplStencil(id, template string, colors map[string]string) error {
	return s.AddTmplStencil(id, template, colors)
}

// AddTableStencil adds a new table stencil
func AddTableStencil(id string, headers, columnOrder []string, colors map[string]string) error {
	return s.AddTableStencil(id, headers, columnOrder, colors)
}
