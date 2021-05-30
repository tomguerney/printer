package setter

import (
	"fmt"
	"io"
	"os"

	"github.com/tomguerney/printer/internal/formatter"
	"github.com/tomguerney/printer/internal/stenciller"
)

// Setter prints formatted and stencilled strings to the set io.Writer
type Setter struct {
	OutWriter  io.Writer
	ErrWriter  io.Writer
	formatter  Formatter
	stenciller Stenciller
}

// New returns a new Setter struct
func New() *Setter {
	return &Setter{
		os.Stdout,
		os.Stderr,
		formatter.New(),
		stenciller.New(),
	}
}

// Formatter formats strings for simple and consistent output
type Formatter interface {
	Msg(string, ...interface{}) string
	Tabulate(headers []string, rows [][]string) []string
	SetTabwriterOptions(twOptions *formatter.TabwriterOptions)
}

// Stenciller formats "data" maps of string key/value pairs according to
// predefined Stencils.
type Stenciller interface {
	AddTmplStencil(id, template string, colors map[string]string) error
	AddTableStencil(id string, headers, columnOrder []string, colors map[string]string) error
	TmplStencil(id string, data map[string]string) (string, error)
	TableStencil(id string, rows []map[string]string) ([][]string, error)
}

// Colorer colors text for terminal output
type Colorer interface {
	Color(text, color string) (string, bool)
}

// SetTabwriterOptions sets tabwriter options
func (s *Setter) SetTabwriterOptions(twOptions *formatter.TabwriterOptions) {
	s.formatter.SetTabwriterOptions(twOptions)
}

// Out prints the passed text appended with a newline to the Writer. If the text
// contains formatting verbs (e.g. %v), they will be formatted as per the
// "...interface{}" variadic parameter in the fashion of fmt.Printf()
func (s *Setter) Out(i interface{}, a ...interface{}) {
	text := fmt.Sprint(i)
	fmt.Fprintf(s.OutWriter, s.formatter.Msg(text, a...))
}

// Err prints the passed text prefixed with "Error: " and appended with a
// newline. If the text contains formatting verbs (e.g. %v), they will be
// formatted as per the "...interface{}" variadic parameter in the fashion of
// fmt.Printf()
func (s *Setter) Err(i interface{}, a ...interface{}) {
	text := fmt.Sprint(i)
	fmt.Fprintf(s.ErrWriter, s.formatter.Msg(text, a...))
}

// Feed prints an empty line to the OutWriter
func (s *Setter) Feed() {
	fmt.Fprintln(s.OutWriter)
}

// Tabulate takes a 2D slice of rows and columns. The 2D slice is tabulated as
// per the tabwriterOptions passed into the domain.Formatter and the internal
// logic of that package. The default tabwriterOptions are set at the root
// printer package level.
//
// Tabulate prints each row from the original 2D slice spaced such that each
// element in each row appear vertically aligned in equally-spaced columns.
func (s *Setter) Tabulate(rows [][]string) {
	tabulated := s.formatter.Tabulate(nil, rows)
	for _, row := range tabulated {
		fmt.Fprintln(s.OutWriter, row)
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
func (s *Setter) TmplStencil(id string, data map[string]string) error {
	result, err := s.stenciller.TmplStencil(id, data)
	if err != nil {
		return err
	}
	fmt.Fprintln(s.OutWriter, result)
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
func (s *Setter) TableStencil(id string, rows []map[string]string) error {
	result, err := s.stenciller.TableStencil(id, rows)
	if err != nil {
		return err
	}
	s.Tabulate(result)
	return nil
}

// AddTmplStencil adds a new Template Stencil with the passed ID and colors.
func (s *Setter) AddTmplStencil(id, template string, colors map[string]string) error {
	return s.stenciller.AddTmplStencil(id, template, colors)
}

// AddTableStencil adds a new table Stencil with the passed ID, headers, and
// colors.
func (s *Setter) AddTableStencil(id string, headers, columnOrder []string, colors map[string]string) error {
	return s.stenciller.AddTableStencil(id, headers, columnOrder, colors)
}
