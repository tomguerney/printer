package printer

import (
	"fmt"
	"io"
	"os"

	"github.com/tomguerney/printer/internal/domain"
	"github.com/tomguerney/printer/internal/formatter"
	"github.com/tomguerney/printer/internal/stenciller"
)

// Printer prints formatted strings to its configured io.Writer output
type Printer struct {
	Out        io.Writer
	formatter  domain.Formatter
	stenciller domain.Stenciller
	logger     
}

// New returns a new Printer struct
func New() *Printer {
	return &Printer{
		os.Stdout,
		formatter.New(),
		stenciller.New(),
	}
}

// Msg prints a formatted message to output
func (p *Printer) Msg(text string, a ...interface{}) {
	fmt.Fprint(p.Out, p.formatter.Msg(text, a))
}

// SMsg returns a formatted message string
func (p *Printer) SMsg(text string, a ...interface{}) string {
	return p.formatter.Msg(text, a)
}

// Error prints a formatted error message to output
func (p *Printer) Error(text string, a ...interface{}) {
	fmt.Fprint(p.Out, p.formatter.Error(text, a))
}

// SError returns a formatted error message string
func (p *Printer) SError(text string, a ...interface{}) string {
	return p.formatter.Error(text, a)
}

// Tabulate takes an array of string arrays and prints a table to output
func (p *Printer) Tabulate(rows [][]string) {
	tabulated := p.formatter.Tabulate(rows)
	for _, row := range tabulated {
		fmt.Fprint(p.Out, row)
	}
}

// STabulate takes an array of string arrays and return an array of formatted rows
func (p *Printer) STabulate(rows [][]string) []string {
	return p.formatter.Tabulate(rows)
}
