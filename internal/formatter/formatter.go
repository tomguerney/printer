package formatter

import (
	"fmt"
	"regexp"
	"strings"
)

// Formatter formats strings for simple and consistent output
type Formatter struct {
	TWOptions *TabwriterOptions
}

// New returns a pointer to a new Formatter struct
func New() *Formatter {
	defaultTabwriterOptions := &TabwriterOptions{
		Minwidth: 0,
		Tabwidth: 8,
		Padding:  4,
		Padchar:  ' ',
		Divchar:  '-',
	}
	return &Formatter{
		defaultTabwriterOptions,
	}
}

// TabwriterOptions opt
type TabwriterOptions struct {
	Minwidth, Tabwidth, Padding int
	Padchar, Divchar            byte
}

// SetTabwriterOptions sets tabwriter options
func (f *Formatter) SetTabwriterOptions(twOptions *TabwriterOptions) {
	f.SetTabwriterOptions(twOptions)
}

// Text returns the passed text appended with a newline. If the text contains
// formatting verbs (e.g. %v), they will be formatted as per the
// "...interface{}" variadic parameter in the fashion of fmt.Printf()
func (f *Formatter) Text(i interface{}, a ...interface{}) string {
	text := fmt.Sprint(i)
	formatted := fmt.Sprintf("%s\n", text)
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
func (f *Formatter) Tabulate(rows [][]string, headers ...string) []string {

	widths := getColWidths(append(rows, headers), f.TWOptions.Minwidth)

	if headers != nil && len(headers) > 0 {
		divRow := createDivRow(widths, f.TWOptions.Minwidth, f.TWOptions.Divchar)
		rows = append([][]string{headers, divRow}, rows...)
	}

	paddedRows := padRows(
		rows,
		widths,
		f.TWOptions.Padding,
		f.TWOptions.Padchar,
	)

	strRows := make([]string, len(rows))

	for i, row := range paddedRows {
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

func padRows(rows [][]string, widths map[int]int, padding int, paddingChar byte) [][]string {
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

func createDivRow(colWidths map[int]int, minWidth int, divChar byte) []string {
	divRow := make([]string, len(colWidths))
	for col, width := range colWidths {
		if width < minWidth {
			divRow[col] = strings.Repeat(string(divChar), minWidth)
		} else {
			divRow[col] = strings.Repeat(string(divChar), width)
		}
	}
	return divRow
}
