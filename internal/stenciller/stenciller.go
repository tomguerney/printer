package stenciller

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/tomguerney/interpolator"
	c "github.com/tomguerney/printer/internal/colorer"
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
// pairs, a "headers" string slice, and a "column order" string slice. When a
// Table Stencil is applied to a slice of "row" maps of string key/value pairs,
// the Stenciller loops through the rows, finding any key in the map that
// matches a key in the Stencil's color map and transforms the data value string
// to the color of the color value. It returns the rows and columns as a 2D
// string slice with a prefixed header row.
type Stenciller struct {
	colorer          colorer
	templateStencils []*TemplateStencil
	tableStencils    []*TableStencil
}

// TemplateStencil is a template stencil
type TemplateStencil struct {
	ID       string
	Template string
	Colors   map[string]string
}

// TableStencil table stencils
type TableStencil struct {
	ID          string
	Colors      map[string]string
	ColumnOrder []string
	Headers     []string
}

type colorer interface {
	Color(text, color string) (string, bool)
}

// New returns a pointer to a new Stenciller struct
func New() *Stenciller {
	return &Stenciller{colorer: c.New()}
}

// Color does a color
func (s *Stenciller) Color(text, color string) (string, bool) {
	return s.colorer.Color(text, color)
}

// AddTemplateStencil adds a new Template Stencil
func (s *Stenciller) AddTemplateStencil(stencil *TemplateStencil) error {
	if stencil.ID == "" {
		return fmt.Errorf("Stencil ID may not be empty")
	}
	for _, existing := range s.templateStencils {
		if stencil.ID == existing.ID {
			return fmt.Errorf("Template Stencil with ID %v already exists", stencil.ID)
		}
	}
	s.templateStencils = append(s.templateStencils, stencil)
	return nil
}

// AddTableStencil adds a new Table Stencil
func (s *Stenciller) AddTableStencil(stencil *TableStencil) error {
	if stencil.ID == "" {
		return fmt.Errorf("Stencil ID may not be empty")
	}
	for _, existing := range s.tableStencils {
		if stencil.ID == existing.ID {
			return fmt.Errorf("Table Stencil with ID %v already exists", stencil.ID)
		}
	}
	s.tableStencils = append(s.tableStencils, stencil)
	return nil
}

// UseTemplateStencil takes the ID of a Template Stencil and a "data" map with string
// key/value pairs. It returns an error if it can't find a Stencil with the
// passed ID or template interpolation fails. It applies the Template Stencil to
// the data map and returns the result.
func (s *Stenciller) UseTemplateStencil(id string, data map[string]string) (string, error) {
	stencil, err := s.findTemplateStencil(id)
	if err != nil {
		return "", err
	}
	coloredData := s.colorMap(stencil.Colors, data)
	interpolated, err := interpolator.Interpolate(stencil.Template, coloredData)
	if err != nil {
		return "", err
	}
	return interpolated, nil
}

// UseTableStencil takes the ID of a Table Stencil and a slice of "row" maps with
// string key/values. It returns an error if it can't find a Stencil with the
// passed ID. It applies the Table Stencil to the row map slice to create a 2D
// slice. If the Headers fields of the Stencil isn't empty, it will will prepend
// the headers to the 2D slice with a dynamically-sized divider row before
// returning the result.
func (s *Stenciller) UseTableStencil(id string, data []map[string]string) (coloredSlices [][]string, err error) {
	stencil, err := s.findTableStencil(id)
	if err != nil {
		return nil, err
	}
	for _, d := range data {
		coloredData := s.colorMap(stencil.Colors, d)
		coloredSlice := mapToSliceInColumnOrder(coloredData, stencil.ColumnOrder)
		coloredSlices = append(coloredSlices, coloredSlice)
	}
	if headerSlices, ok := createHeaderSlices(stencil, data); ok {
		coloredSlices = append(headerSlices, coloredSlices...)
	}
	return coloredSlices, nil

}

func (s *Stenciller) findTemplateStencil(id string) (*TemplateStencil, error) {
	for _, stencil := range s.templateStencils {
		if stencil.ID == id {
			return stencil, nil
		}
	}
	return nil, fmt.Errorf("Unable to find template stencil with id of %v", id)
}

func (s *Stenciller) findTableStencil(id string) (*TableStencil, error) {
	for _, stencil := range s.tableStencils {
		if stencil.ID == id {
			return stencil, nil
		}
	}
	return nil, fmt.Errorf("Unable to find table stencil with id of %v", id)
}

func (s *Stenciller) colorMap(colors map[string]string, data map[string]string) map[string]string {
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

func mapToSliceInColumnOrder(mapRow map[string]string, columnOrder []string) []string {
	sliceRow := make([]string, len(columnOrder))
	for key, value := range mapRow {
		col, err := indexOf(key, columnOrder)
		if err != nil {
			log.Debug().Err(err)
			continue
		}
		sliceRow[col] = value
	}
	return sliceRow
}

func createHeaderSlices(stencil *TableStencil, dataMaps []map[string]string) (_ [][]string, ok bool) {
	if len(stencil.Headers) == 0 {
		return nil, false
	}
	dataSlices := [][]string{}
	for _, m := range dataMaps {
		dataSlices = append(dataSlices, mapToSliceInColumnOrder(m, stencil.ColumnOrder))
	}
	dataSlices = append([][]string{stencil.Headers}, dataSlices...)
	colWidths := getColWidths(dataSlices)
	divRow := createDivRow(colWidths)
	return [][]string{stencil.Headers, divRow}, true
}


func createDivRow(colWidths map[int]int) []string {
	divRow := make([]string, len(colWidths))
	for col, width := range colWidths {
		divRow[col] = strings.Repeat("-", width)
	}
	return divRow
}

func getColWidths(rows [][]string) map[int]int {
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

func indexOf(elem string, data []string) (int, error) {
	for k, v := range data {
		if elem == v {
			return k, nil
		}
	}
	return 0, fmt.Errorf("Unable to find index of element %v", elem)
}
