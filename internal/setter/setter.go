package setter

import (
	"fmt"
	"io"
	"os"

	"github.com/tomguerney/printer/internal/domain"
	"github.com/tomguerney/printer/internal/formatter"
	"github.com/tomguerney/printer/internal/stenciller"
)

// Setter prints formatted strings to its configured io.Writer output
type Setter struct {
	Writer     io.Writer
	formatter  domain.Formatter
	stenciller domain.Stenciller
}

// New returns a new Setter struct
func New() *Setter {
	return &Setter{
		os.Stdout,
		formatter.New(),
		stenciller.New(),
	}
}

// Msg prints a formatted message to output
func (s *Setter) Msg(text string, a ...interface{}) {
	fmt.Fprintf(s.Writer, s.formatter.Msg(text, a...))
}

// SMsg returns a formatted message string
func (s *Setter) SMsg(text string, a ...interface{}) string {
	return s.formatter.Msg(text, a...)
}

// Error prints a formatted error message to output
func (s *Setter) Error(text string, a ...interface{}) {
	fmt.Fprintf(s.Writer, s.formatter.Error(text, a...))
}

// SError returns a formatted error message string
func (s *Setter) SError(text string, a ...interface{}) string {
	return s.formatter.Error(text, a...)
}

// Tabulate takes an array of string arrays and prints a table to output
func (s *Setter) Tabulate(rows [][]string) {
	tabulated := s.formatter.Tabulate(rows)
	for _, row := range tabulated {
		fmt.Fprint(s.Writer, row)
	}
}

// CTabulate takes an array of string arrays, colors them, and prints a table to output
func (s *Setter) CTabulate(rows [][]string, colors map[string]string) {

	tabulated := s.formatter.Tabulate(rows)
	for _, row := range tabulated {
		fmt.Fprint(s.Writer, row)
	}
}

// STabulate takes an array of string arrays and return an array of formatted rows
func (s *Setter) STabulate(rows [][]string) []string {
	return s.formatter.Tabulate(rows)
}

// AddTmplStencil adds a new template stencil
func (s *Setter) AddTmplStencil(id, template string, colors map[string]string) error {
	return s.stenciller.AddTmplStencil(id, template, colors)
}

// AddTableStencil adds a new table stencil
func (s *Setter) AddTableStencil(id string, colors map[string]string) error {
	return s.stenciller.AddTableStencil(id, colors)
}

// Stencil applies a string map to the stencil with the passed ID and prints it to output
func (s *Setter) Stencil(id string, data map[string]string) error {
	result, err := s.stenciller.TmplStencil(id, data)
	if err != nil {
		return err
	}
	fmt.Fprint(s.Writer, result)
	return nil
}

// FStencil applies a string map to the stencil with the passed ID and returns the result
func (s *Setter) FStencil(id string, data map[string]string) (string, error) {
	result, err := s.stenciller.TmplStencil(id, data)
	if err != nil {
		return "", err
	}
	return result, nil
}

func AddTableStencil(id string, colors map[string]string) {

}

// StencilTable take an array of string maps and prints stencilled rows to output
func StencilTable(id string, rows []map[string]string, colors map[string]string) {

}
