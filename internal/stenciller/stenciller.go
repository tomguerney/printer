package stenciller

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"

	"github.com/tomguerney/printer/internal/colorer"
	"github.com/tomguerney/printer/internal/domain"
)

// Stenciller formats "data" maps of string key/value pairs according to
// predefined Stencils.
//
// Stencils are either Template Stencils or Table Stencils.
//
// A Template Stencil is comprised of an ID, a "color" map of string key/value
// pairs, and a template string as per the "text/template" package from the Go
// standard library. When a Template Stencil is applied to a "data" map of
// string key/value pairs, it finds any key in the map that matches a key in the
// Template Stencil's color map and transforms the data value string to the
// color of the color value. The data map is then applied to the template to
// produce a single string.
//
// A Table Stencil is comprised of an ID, a "color" map of string key/value
// pairs, and a "headers" string slice. When a Table Stencil is applied to a
// slice of "row" maps of string key/value pairs, the Stenciller loops through
// the rows, finding any key in the map that matches a key in the Stencil's
// color map and transforms the data value string to the color of the color
// value. It returns the rows and columns as a 2D string slice with a prefixed
// header row.
type Stenciller struct {
	tmplStencils  []*tmplStencil
	tableStencils []*tableStencil
	colorer       domain.Colorer
}

type tmplStencil struct {
	ID       string
	Template string
	Colors   map[string]string
}

type tableStencil struct {
	ID          string
	Colors      map[string]string
	ColumnOrder []string
	Headers     []string
}

// New returns a pointer to a new Stenciller struct
func New() *Stenciller {
	return &Stenciller{
		colorer: colorer.New(),
	}
}

// AddTmplStencil adds a new Template Stencil
func (s *Stenciller) AddTmplStencil(
	id, template string,
	colors map[string]string,
) error {
	if id == "" {
		return fmt.Errorf("Stencil ID may not be empty")
	}
	for _, stencil := range s.tmplStencils {
		if stencil.ID == id {
			return fmt.Errorf("Template Stencil with ID %v already exists", id)
		}
	}
	s.tmplStencils = append(s.tmplStencils, &tmplStencil{id, template, colors})
	return nil
}

// AddTableStencil adds a new Table Stencil
func (s *Stenciller) AddTableStencil(
	id string,
	headers []string,
	columnOrder []string,
	colors map[string]string,
) error {
	if id == "" {
		return fmt.Errorf("Stencil ID may not be empty")
	}
	for _, stencil := range s.tableStencils {
		if stencil.ID == id {
			return fmt.Errorf("Table Stencil with ID %v already exists", id)
		}
	}
	s.tableStencils = append(
		s.tableStencils,
		&tableStencil{
			ID:          id,
			Headers:     headers,
			Colors:      colors,
			ColumnOrder: columnOrder,
		})
	return nil
}

// TmplStencil takes the ID of a Template Stencil and a "data" map with string
// key/value pairs. It returns an error if it can't find a Stencil with the
// passed ID or template interpolation fails. It applies the Template Stencil to
// the data map and returns the result.
func (s *Stenciller) TmplStencil(
	id string,
	data map[string]string,
) (string, error) {
	stencil, err := s.findTmplStencil(id)
	if err != nil {
		return "", err
	}
	coloredData := s.colorData(stencil.Colors, data)
	interpolated, err := s.interpolate(stencil.Template, coloredData)
	if err != nil {
		return "", err
	}
	return interpolated, nil
}

// TableStencil takes the ID of a Table Stencil and a slice of "row" maps with
// string key/values. It returns an error if it can't find a Stencil with the
// passed ID. It applies the Table Stencil to the row map slice to create a 2D
// slice. If the Headers fields of the Stencil isn't empty, it will will prepend
// the headers to the 2D slice with a dynamically-sized divider row before
// returning the result.
func (s *Stenciller) TableStencil(
	id string,
	rawDataRows []map[string]string,
) (colorSliceRows [][]string, err error) {
	stencil, err := s.findTableStencil(id)
	if err != nil {
		return nil, err
	}
	for _, rawDataRow := range rawDataRows {
		colorDataRow := s.colorData(stencil.Colors, rawDataRow)
		colorSliceRow := make([]string, len(stencil.ColumnOrder))
		for key, value := range colorDataRow {
			col, err := s.indexOf(key, stencil.ColumnOrder)
			if err != nil {
				log.Info().Err(err)
				continue
			}
			colorSliceRow[col] = value
		}
		colorSliceRows = append(colorSliceRows, colorSliceRow)
	}
	if len(stencil.Headers) > 0 {
		colorSliceRows = s.prependHeaders(colorSliceRows, stencil.Headers)
	}
	return colorSliceRows, nil
}

func (s *Stenciller) prependHeaders(
	rows [][]string,
	headers []string,
) [][]string {
	rowsWithHeader := [][]string{headers, s.createDivRow(append(rows, headers))}
	return append(rowsWithHeader, rows...)
}

func (s *Stenciller) findTmplStencil(id string) (*tmplStencil, error) {
	for _, stencil := range s.tmplStencils {
		if stencil.ID == id {
			return stencil, nil
		}
	}
	return nil, fmt.Errorf("Unable to find template stencil with id of %v", id)
}

func (s *Stenciller) findTableStencil(id string) (*tableStencil, error) {
	for _, stencil := range s.tableStencils {
		if stencil.ID == id {
			return stencil, nil
		}
	}
	return nil, fmt.Errorf("Unable to find table stencil with id of %v", id)
}

func (s *Stenciller) colorData(
	colors map[string]string,
	data map[string]string,
) map[string]string {
	colored := make(map[string]string, len(data))
	for key, val := range data {
		if col, ok := colors[key]; ok {
			colored[key] = s.colorValue(val, col)
		} else {
			colored[key] = val
		}
	}
	return colored
}

func (s *Stenciller) colorValue(val, col string) string {
	if coloredVal, ok := s.colorer.Color(val, col); ok {
		return coloredVal
	}
	log.Info().Msgf("Unable to set [value=%v] with color [%v]", val, col)
	return val
}

func (s *Stenciller) interpolate(
	tmpl string,
	data map[string]string,
) (string, error) {
	builder := strings.Builder{}
	parsed, err := template.New("stencil").Parse(tmpl)
	if err != nil {
		return "", err
	}
	err = parsed.Execute(&builder, data)
	if err != nil {
		return "", err
	}
	return builder.String(), nil
}

func (s *Stenciller) createDivRow(rows [][]string) []string {
	colWidths := s.getColWidths(rows)
	divRow := make([]string, len(colWidths))
	for col, width := range colWidths {
		divRow[col] = strings.Repeat("-", width)
	}
	return divRow
}

func (s *Stenciller) getColWidths(
	rows [][]string,
) map[int]int {
	maxCols := 0
	for _, row := range rows {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}
	widths := make(map[int]int, maxCols)
	for _, row := range rows {
		for col, elem := range row {
			if len(elem) > widths[col] {
				widths[col] = len(elem)
			}
		}
	}
	return widths
}

func (s *Stenciller) indexOf(elem string, data []string) (int, error) {
	for k, v := range data {
		if elem == v {
			return k, nil
		}
	}
	return 0, fmt.Errorf("Unable to find index of element %v", elem)
}
