package printer

import (
	"fmt"
	"io"
	"os"

	"github.com/tomguerney/printer/internal/formatter"
	"github.com/tomguerney/printer/internal/stenciller"
)

// Printer prints formatted and stencilled strings to the set io.Writer
type Printer struct {
	OutWriter  io.Writer
	ErrWriter  io.Writer
	formatter  Formatter
	stenciller Stenciller
}

// Formatter formats strings for simple and consistent output
type Formatter interface {
	Text(interface{}, ...interface{}) string
	Tabulate(rows [][]string, headers ...string) []string
	SetTabwriterOptions(twOptions *formatter.TabwriterOptions)
}

// Stenciller formats "data" maps of string key/value pairs according to
// predefined Stencils.
type Stenciller interface {
	AddTmplStencil(id, template string, colors map[string]string) error
	AddTableStencil(id string, headers, columnOrder []string, colors map[string]string) error
	TmplStencil(id string, data map[string]string) (string, error)
	TableStencil(id string, rows []map[string]string) ([][]string, error)
	Color(text, color string) (string, bool)
}

// New a new printer
func New() *Printer {
	return &Printer{
		os.Stdout,
		os.Stderr,
		formatter.New(),
		stenciller.New(),
	}
}

// GetPrinter returns the singleton Prtiner
func GetPrinter() *Printer {
	return singleton
}

var singleton = New()

// SetOutWriter sets the OutWriter
func SetOutWriter(writer io.Writer) {
	singleton.SetOutWriter(writer)
}

func (p *Printer) SetOutWriter(writer io.Writer) {
	p.OutWriter = writer
}

// SetErrWriter sets the ErrWriter
func SetErrWriter(writer io.Writer) {
	singleton.SetErrWriter(writer)
}

func (p *Printer) SetErrWriter(writer io.Writer) {
	p.OutWriter = writer
}

// SetTabwriterOptions sets tabwriter options
func SetTabwriterOptions(twOptions *formatter.TabwriterOptions) {
	singleton.formatter.SetTabwriterOptions(twOptions)
}

func (p *Printer) SetTabwriterOptions(twOptions *formatter.TabwriterOptions) {
	p.formatter.SetTabwriterOptions(twOptions)
}

// Out prints the passed text appended with a newline to the Writer. If the text
// contains formatting verbs (e.g. %v), they will be formatted as per the
// "...interface{}" variadic parameter in the fashion of fmt.Printf()
func Out(i interface{}, a ...interface{}) {
	singleton.Out(i, a...)
}

func (p *Printer) Out(i interface{}, a ...interface{}) {
	fmt.Fprintf(p.OutWriter, p.formatter.Text(i, a...))
}

// Err prints the passed text prefixed with "Error: " and appended with a
// newline. If the text contains formatting verbs (e.g. %v), they will be
// formatted as per the "...interface{}" variadic parameter in the fashion of
// fmt.Printf()
func Err(i interface{}, a ...interface{}) {
	singleton.Err(i, a...)
}

func (p *Printer) Err(i interface{}, a ...interface{}) {
	fmt.Fprintf(p.ErrWriter, p.formatter.Text(i, a...))
}

// Feed prints an empty line to the OutWriter
func Feed() {
	singleton.Feed()
}

func (p *Printer) Feed() {
	fmt.Fprintln(p.OutWriter)
}

// Will return plain text if color not found
func Color(text, color string) string {
	return singleton.Color(text, color)
}

func (p *Printer) Color(text, color string) string {
	colorized, _ := p.stenciller.Color(text, color)
	return colorized
}

// Tabulate takes a 2D slice of rows and columns. The 2D slice is tabulated as
// per the tabwriterOptions passed into the domain.Formatter and the internal
// logic of that package. The default tabwriterOptions are set at the root
// printer package level.
//
// Tabulate prints each row from the original 2D slice spaced such that each
// element in each row appear vertically aligned in equally-spaced columns.
func Tabulate(rows [][]string, headers ...string) {
	singleton.Tabulate(rows, headers...)
}

func (p *Printer) Tabulate(rows [][]string, headers ...string) {
	tabulated := p.formatter.Tabulate(rows, headers...)
	for _, row := range tabulated {
		fmt.Fprintln(p.OutWriter, row)
	}
}

// TmplStencil takes the ID of a Template Stencil and a "data" map with string
// key/value pairs. It returns an error if it can't find a Stencil with the
// passed ID. It applies the Template Stencil to the map and prints the result.
//
// A Template Stencil is comprised of an ID, a "color" map of string key/value
// pairs, and a template string as per the "text/template" package from the Go
// standard library. When a Template Stencil is applied to a data map, it finds
// any key in the map that matches a key in the Template Stencil's color map and
// transforms the data value string to the color of the color value. The data
// map is then applied to the template to produce a single string.
func TmplStencil(id string, data map[string]string) error {
	return singleton.TmplStencil(id, data)
}

func (p *Printer) TmplStencil(id string, data map[string]string) error {
	result, err := p.stenciller.TmplStencil(id, data)
	if err != nil {
		return err
	}
	fmt.Fprintln(p.OutWriter, result)
	return nil
}

// TableStencil takes the ID of a Table Stencil and a slice of "row" maps with
// string key/values. It returns an error if it can't find a Stencil with the
// passed ID. It applies the Table Stencil to the map slice and tabulates and
// prints the result.
//
// A Table Stencil is comprised of an ID, a "color" map of string key/values,
// and a "headers" string slice. When a Table Stencil is applied to a slice of
// row maps, the Stenciller loops through the rows, finding any key in the map
// that matches a key in the Stencil's color map and transforms the data value
// string to the color of the color value. It returns the rows and columns as a
// 2D string slice with a prefixed header row.
func TableStencil(id string, rows []map[string]string) error {
	return singleton.TableStencil(id, rows)
}

func (p *Printer) TableStencil(id string, rows []map[string]string) error {
	result, err := p.stenciller.TableStencil(id, rows)
	if err != nil {
		return err
	}
	p.Tabulate(result)
	return nil
}

// AddTmplStencil adds a new Template Stencil with the passed ID and colors.
func AddTmplStencil(id, template string, colors map[string]string) error {
	return singleton.AddTmplStencil(id, template, colors)
}

func (p *Printer) AddTmplStencil(id, template string, colors map[string]string) error {
	return p.stenciller.AddTmplStencil(id, template, colors)
}

// AddTableStencil adds a new table Stencil with the passed ID, headers, and
// colors.
func AddTableStencil(id string, headers, columnOrder []string, colors map[string]string) error {
	return singleton.AddTableStencil(id, headers, columnOrder, colors)
}

func (p *Printer) AddTableStencil(id string, headers, columnOrder []string, colors map[string]string) error {
	return p.stenciller.AddTableStencil(id, headers, columnOrder, colors)
}
