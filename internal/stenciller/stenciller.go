package stenciller

import (
	"fmt"
	"html/template"
	"strings"

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
// pairs, and a template string as per the "html/template" package from the Go
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
	ID      string
	Colors  map[string]string
	Headers []string
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
			return fmt.Errorf("Stencil with ID %v already exists", id)
		}
	}
	s.tmplStencils = append(s.tmplStencils, &tmplStencil{id, template, colors})
	return nil
}

// AddTableStencil adds a new Table Stencil
func (s *Stenciller) AddTableStencil(
	id string,
	headers []string,
	colors map[string]string,
) error {
	if id == "" {
		return fmt.Errorf("Stencil ID may not be empty")
	}
	for _, stencil := range s.tableStencils {
		if stencil.ID == id {
			return fmt.Errorf("Stencil with ID %v already exists", id)
		}
	}
	s.tableStencils = append(
		s.tableStencils,
		&tableStencil{ID: id, Headers: headers, Colors: colors})
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
	if stencil.Template == "" {
		log.Info().Msgf("Template with ID %v is an empty string", id)
	}
	result, err := s.interpolate(stencil.Template, coloredData)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v\n", result), nil
}

// TableStencil takes the ID of a Table Stencil and a slice of "row" maps with
// string key/values. It returns an error if it can't find a Stencil with the
// passed ID. It applies the Table Stencil to the row map slice and tabulates
// and returns the result.
func (s *Stenciller) TableStencil(
	id string,
	dataRows []map[string]string,
) ([][]string, error) {
	stencil, err := s.findTableStencil(id)
	if err != nil {
		return nil, err
	}
	coloredRows := [][]string{}
	for _, row := range dataRows {
		coloredData := s.colorData(stencil.Colors, row)
		coloredRow := make([]string, len(coloredData))
		for _, value := range coloredData {
			coloredRow = append(coloredRow, value)
		}
		coloredRows = append(coloredRows, coloredRow)
	}
	return coloredRows, nil
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
