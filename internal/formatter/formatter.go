package formatter

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/tomguerney/printer/internal/domain"
)

// Formatter formats strings for simple and consistent output
type Formatter struct {
	tabwriterOptions *domain.TabwriterOptions
}

// New returns a pointer to a new Formatter struct
func New(tabwriterOptions *domain.TabwriterOptions) *Formatter {
	return &Formatter{
		tabwriterOptions,
	}
}

// Msg returns a string comprised of the passed text ending with a newline
// The text may contain formatting verbs that will be formatted with
// the passed operands
func (f *Formatter) Msg(text string, a ...interface{}) string {
	formatted := fmt.Sprintf("%s\n", text)
	return fmt.Sprintf(formatted, a...)
}

// Error returns a string comprised of the passed text prefixed with "Error: " with a newline
// The text may contain formatting verbs that will be formatted with
// the passed operands
func (f *Formatter) Error(text string, a ...interface{}) string {
	formatted := fmt.Sprintf("Error: %s\n", text)
	return fmt.Sprintf(formatted, a...)
}

// Tabulate returns a formmatted table
func (f *Formatter) Tabulate(rows [][]string) []string {

	builder := strings.Builder{}
	writer := tabwriter.NewWriter(
		&builder,
		f.tabwriterOptions.Minwidth,
		f.tabwriterOptions.Tabwidth,
		f.tabwriterOptions.Padding,
		f.tabwriterOptions.Padchar,
		f.tabwriterOptions.Flags,
	)

	for _, row := range rows {
		fmt.Fprintln(writer, fmt.Sprintf("%s", strings.Join(row, "\t")))
	}

	writer.Flush()
	table := strings.Split(builder.String(), "\n")

	return table[:len(table)-1]
}
