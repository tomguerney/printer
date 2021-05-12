package printer

import (
	"io"
	"os"

	"github.com/tomguerney/printer/internal/setter"
)

var s = setter.New()

// Writer is the io.Writer to print to
var Writer io.Writer = os.Stdout

// SetWriter sets the writer of the setter
func SetWriter(writer io.Writer) {
	s.Writer = writer
}

// Msg prints a formatted message to output
func Msg(text string, a ...interface{}) {
	s.Msg(text, a...)
}

// SMsg returns a formatted message string
func SMsg(text string, a ...interface{}) string {
	return s.SMsg(text, a...)
}

// Error prints a formatted error message to output
func Error(text string, a ...interface{}) {
	s.Error(text, a...)
}

// SError returns a formatted error message string
func SError(text string, a ...interface{}) string {
	return s.SError(text, a...)
}

// Tabulate takes an array of string arrays and prints a table to output
func Tabulate(rows [][]string) {
	s.Tabulate(rows)
}

// STabulate takes an array of string arrays and return an array of formatted rows
func STabulate(rows [][]string) []string {
	return s.STabulate(rows)
}

// AddTmplStencil adds a new stencil
func AddTmplStencil(id, template string, colors map[string]string) error {
	return s.AddTmplStencil(id, template, colors)
}

// TmplStencil applies a string map to the stencil with the passed ID and prints it to output
func TmplStencil(id string, data map[string]string) error {
	return s.TmplStencil(id, data)
}

// STmplStencil applies a string map to the stencil with the passed ID and returns the result
func STmplStencil(id string, data map[string]string) (string, error) {
	return s.STmplStencil(id, data)
}

// TableStencil take an array of string maps and prints stencilled rows to output
func TableStencil(id string, rows []map[string]string) {
	// return s.StencilTable(id, rows)
}
