package printer

import (
	"io"
	"os"

	s "github.com/tomguerney/printer/internal/setter"
)

var setter = s.New()

// Writer is the io.Writer to print to
var Writer io.Writer = os.Stdout

// SetWriter sets the writer of the setter
func SetWriter(writer io.Writer) {
	setter.Writer = writer
}

// Msg prints a formatted message to output
func Msg(text string, a ...interface{}) {
	setter.Msg(text, a...)
}

// SMsg returns a formatted message string
func SMsg(text string, a ...interface{}) string {
	return setter.SMsg(text, a...)
}

// Error prints a formatted error message to output
func Error(text string, a ...interface{}) {
	setter.Error(text, a...)
}

// SError returns a formatted error message string
func SError(text string, a ...interface{}) string {
	return setter.SError(text, a...)
}

// Tabulate takes an array of string arrays and prints a table to output
func Tabulate(rows [][]string) {
	setter.Tabulate(rows)
}

// STabulate takes an array of string arrays and return an array of formatted rows
func STabulate(rows [][]string) []string {
	return setter.STabulate(rows)
}

// AddStencil adds a new stencil
func AddStencil(id, template string, colors map[string]string) error {
	return setter.AddStencil(id, template, colors)
}

// UseStencil applies a string map to the stencil with the passed ID and prints it to output
func UseStencil(id string, data map[string]string) error {
	return setter.UseStencil(id, data)
}

// FUseStencil applies a string map to the stencil with the passed ID and returns the result
func FUseStencil(id string, data map[string]string) (string, error) {
	return setter.FUseStencil(id, data)
}
