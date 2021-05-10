package printer

import (
	"fmt"
	"io"
	"os"

	f "github.com/tomguerney/printer/internal/formatter"
	p "github.com/tomguerney/printer/internal/printer"
	s "github.com/tomguerney/printer/internal/stenciller"
)

var formatter = f.New()
var stenciller = s.New()
var printer = p.New()

// Writer is the io.Writer to print to
var Writer io.Writer = os.Stdout

// Msg prints a formatted message to output
func Msg(text string, a ...interface{}) {
	fmt.Fprint(Writer, formatter.Msg(text, a))
}

// SMsg returns a formatted message string
func SMsg(text string, a ...interface{}) string {
	return formatter.Msg(text, a)
}

// Error prints a formatted error message to output
func Error(text string, a ...interface{}) {
	fmt.Fprint(Writer, formatter.Error(text, a))
}

// SError returns a formatted error message string
func SError(text string, a ...interface{}) string {
	return formatter.Error(text, a)
}

// Tabulate takes an array of string arrays and prints a table to output
func Tabulate(rows [][]string) {
	tabulated := formatter.Tabulate(rows)
	for _, row := range tabulated {
		fmt.Fprint(Writer, row)
	}
}

// STabulate takes an array of string arrays and return an array of formatted rows
func STabulate(rows [][]string) []string {
	return formatter.Tabulate(rows)
}

// AddStencil adds a new stencil
func AddStencil(id, template string, colors map[string]string) error {
	return stenciller.AddStencil(id, template, colors)
}

// UseStencil applies a string map to the stencil with the passed ID and prints it to output
func UseStencil(id string, data map[string]string) error {
	result, err := stenciller.UseStencil(id, data)
	if err != nil {
		return err
	}
	fmt.Fprint(Writer, result)
	return nil
}

// FUseStencil applies a string map to the stencil with the passed ID and returns the result
func FUseStencil(id string, data map[string]string) (string, error) {
	result, err := stenciller.UseStencil(id, data)
	if err != nil {
		return "", err
	}
	return result, nil
}
