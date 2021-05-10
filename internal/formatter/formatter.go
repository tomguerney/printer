package formatter

import (
	"fmt"
	"strings"
	"text/tabwriter"
)

// Formatter formats text for output
type Formatter struct{}

// New returns a pointer to a new Formatter struct
func New() *Formatter {
	return &Formatter{}
}

// Msg creates a formatted message
func (f *Formatter) Msg(text string, a ...interface{}) string {
	formatted := fmt.Sprintf("%s\n", text)
	return fmt.Sprintf(formatted, a...)
}

// Error creates a formatted message prefixed with "Error: "
func (f *Formatter) Error(text string, a ...interface{}) string {
	formatted := fmt.Sprintf("Error: %s\n", text)
	return fmt.Sprintf(formatted, a...)
}

// Tabulate creates a formmatted table
func (f *Formatter) Tabulate(rows [][]string) []string {

	builder := strings.Builder{}
	writer := tabwriter.NewWriter(&builder, 0, 8, 4, ' ', 0)

	for _, row := range rows {
		fmt.Fprintln(writer, fmt.Sprintf("%s", strings.Join(row, "\t")))
	}

	writer.Flush()
	table := strings.Split(builder.String(), "\n")

	return table[:len(table)-1]
}
