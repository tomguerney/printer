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

// Msg returns the passed text appended with a newline. If the text contains
// formatting verbs (e.g. %v), they will be formatted as per the
// "...interface{}" variadic parameter in the fashion of fmt.Printf()
func (f *Formatter) Msg(text string, a ...interface{}) string {
	formatted := fmt.Sprintf("%s\n", text)
	return fmt.Sprintf(formatted, a...)
}

// Error returns the passed text prefixed with "Error: " and appended with a
// newline. If the text contains formatting verbs (e.g. %v), they will be
// formatted as per the "...interface{}" variadic parameter in the fashion of
// fmt.Printf()
func (f *Formatter) Error(text string, a ...interface{}) string {
	formatted := fmt.Sprintf("Error: %s\n", text)
	return fmt.Sprintf(formatted, a...)
}

// Tabulate takes a 2D slice of rows and columns. The 2D slice is tabulated as
// per the tabwriterOptions passed into the NewWriter function from the
// "text/tabwriter" package from the Go standard library. The default
// tabwriterOptions are set at the root printer package level.
//
// Tabulate returns a one-dimensional slice of strings, with each element formed
// from a row of strings from the original 2D slice. Each row is spaced such
// that when the slice is printed row by row, the element in each row appear
// vertically aligned in equally-spaced columns
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

