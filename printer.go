package printer

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

// Printer is a tool for formatting output
type Printer struct {
	out io.Writer
}

// New returns a newly created Printer object.
func New() *Printer {
	p := &Printer{os.Stdout}
	return p
}

// Msg prints a formatted message
func (p *Printer) Msg(text string, a ...interface{}) {
	textln := fmt.Sprintf("%s\n", text)
	fmt.Fprintf(p.out, textln, a...)
}

// Error prints a formatted message prefixed with "Error: "
func (p *Printer) Error(text string, a ...interface{}) {
	textln := fmt.Sprintf("Error: %s\n", text)
	fmt.Fprintf(p.out, textln, a...)
}

// Tabulate prints a formmatted table
func (p *Printer) Tabulate(rows [][]string) ([]string, error) {

	builder := strings.Builder{}
	writer := tabwriter.NewWriter(&builder, 0, 8, 4, ' ', 0)

	for _, row := range rows {
		fmt.Fprintln(writer, fmt.Sprintf("%s", strings.Join(row, "\t")))
	}

	writer.Flush()
	table := strings.Split(builder.String(), "\n")

	return table[:len(table)-1], nil
}
