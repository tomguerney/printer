package printer

import (
	"fmt"
	"io"

	"github.com/tomguerney/printer/internal/domain"
	"github.com/tomguerney/printer/internal/formatter"
	"github.com/tomguerney/printer/internal/stenciller"
)

// Printer prints formatted strings to its configured io.Writer output
type Printer struct {
	out        io.Writer
	formatter  domain.Formatter
	stenciller domain.Stenciller
}

// New returns a new Printer struct
func New(out io.Writer) *Printer {
	return &Printer{
		out: out,
		formatter: formatter.New(),
		stenciller: stenciller.New(),
	}
}

// Msg prints a formatted message to output
func (p *Printer) Msg(text string, a ...interface{}) {
	fmt.Fprint(p.out, p.formatter.Msg(text, a))
}

// SMsg returns a formatted message string
func (p *Printer) SMsg(text string, a ...interface{}) string {
	return p.formatter.Msg(text, a)
}

// Error prints a formatted error message to output
func (p *Printer) Error(text string, a ...interface{}) {
	fmt.Fprint(p.out, p.formatter.Error(text, a))
}

// SError returns a formatted error message string
func (p *Printer) SError(text string, a ...interface{}) string {
	return p.formatter.Error(text, a)
}

// Tabulate takes an array of string arrays and prints a table to output
func (p *Printer) Tabulate(rows [][]string) {
	tabulated := p.formatter.Tabulate(rows)
	for _, row := range tabulated {
		fmt.Fprint(p.out, row)
	}
}

// STabulate takes an array of string arrays and return an array of formatted rows
func (p *Printer) STabulate(rows [][]string) []string {
	return p.formatter.Tabulate(rows)
}

// AddStencil adds a new stencil
func (p *Printer) AddStencil(id, template string, colors map[string]string) error {
	return p.stenciller.AddStencil(id, template, colors)
}

// UseStencil applies a string map to the stencil with the passed ID and prints it to output
func (p *Printer) UseStencil(id string, data map[string]string) error {
	result, err := p.stenciller.UseStencil(id, data)
	if err != nil {
		return err
	}
	fmt.Fprint(p.out, result)
	return nil
}

// FUseStencil applies a string map to the stencil with the passed ID and returns the result
func (p *Printer) FUseStencil(id string, data map[string]string) (string, error) {
	result, err := p.stenciller.UseStencil(id, data)
	if err != nil {
		return "", err
	}
	return result, nil
}
