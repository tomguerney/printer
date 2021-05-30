package formatter

import (
	"fmt"
	"regexp"
	"strings"

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
	if len(a) > 0 {
		return fmt.Sprintf(formatted, a...)
	} else {
		return fmt.Sprint(formatted)
	}

}

// Error returns the passed text prefixed with "Error: " and appended with a
// newline. If the text contains formatting verbs (e.g. %v), they will be
// formatted as per the "...interface{}" variadic parameter in the fashion of
// fmt.Printf()
func (f *Formatter) Error(text string, a ...interface{}) string {
	formatted := fmt.Sprintf("Error: %s\n", text)
	if len(a) > 0 {
		return fmt.Sprintf(formatted, a...)
	} else {
		return fmt.Sprint(formatted)
	}
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
func (f *Formatter) Tabulate(headers []string, rows [][]string) []string {

	widths := getColWidths(append(rows, headers), f.tabwriterOptions.Minwidth)

	if headers != nil && len(headers) < 0 {
		divRow := createDivRow(widths, f.tabwriterOptions.Minwidth)
		headerRows := [][]string{headers, divRow}
		rows = append(headerRows, rows...)
	}

	spacedRows := spaceCols(
		rows,
		widths,
		f.tabwriterOptions.Padding,
		f.tabwriterOptions.Padchar,
	)

	strRows := make([]string, len(rows))

	for i, row := range spacedRows {
		strRows[i] = strings.Join(row, "")
	}
	return strRows
}

func getColWidths(rows [][]string, minWidth int) map[int]int {
	widths := make(map[int]int)
	for _, row := range rows {
		for col, elem := range row {
			if _, ok := widths[col]; !ok {
				widths[col] = minWidth
			}
			if lenNoAnsi(elem) > widths[col] {
				widths[col] = lenNoAnsi(elem)
			}
		}
	}
	return widths
}

func lenNoAnsi(str string) int {
	const ansi = "[\u001b\u009b][[()#;?]*(?:[0-9]{1,4}(?:;[0-9]{0,4})*)?[0-9A-ORZcf-nqry=><]"
	var re = regexp.MustCompile(ansi)
	return len(re.ReplaceAllString(str, ""))
}

func spaceCols(rows [][]string, widths map[int]int, padding int, paddingChar byte) [][]string {
	for _, row := range rows {
		for col, val := range row {
			diff := 0
			if l := lenNoAnsi(val); l < widths[col] {
				diff = widths[col] - l
			}
			row[col] = val
			if col < len(row)-1 {
				row[col] = row[col] + strings.Repeat(string(paddingChar), diff+padding)
			}
		}
	}
	return rows
}

func createDivRow(colWidths map[int]int, minWidth int) []string {
	divRow := make([]string, len(colWidths))
	for col, width := range colWidths {
		if width < minWidth {
			divRow[col] = strings.Repeat("-", minWidth)
		} else {
			divRow[col] = strings.Repeat("-", width)
		}
	}
	return divRow
}
