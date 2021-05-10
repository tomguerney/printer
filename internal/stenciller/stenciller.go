package stenciller

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/apex/log"
	"github.com/tomguerney/printer/internal/domain"
)

// Stenciller takes structs and formats their field values to predefined templates and colors
type Stenciller struct {
	stencils []*stencil
	colorer  domain.Colorer
	logger   log.Interface
}

type stencil struct {
	ID       string
	Template string
	Colors   map[string]string
}

// New returns a pointer to a new Stenciller struct
func New() *Stenciller {
	return &Stenciller{}
}

// AddStencil adds a new stencil
func (s *Stenciller) AddStencil(id, template string, colors map[string]string) error {
	if id == "" {
		return fmt.Errorf("Stencil ID may not be empty")
	}
	for _, stencil := range s.stencils {
		if stencil.ID == id {
			return fmt.Errorf("Stencil with ID %v already exists", id)
		}
	}
	s.stencils = append(s.stencils, &stencil{id, template, colors})
	return nil
}

// UseStencil applies a string map to the stencil with the passed ID
func (s *Stenciller) UseStencil(id string, data map[string]string) (string, error) {
	stencil, err := s.findStencil(id)
	if err != nil {
		return "", err
	}
	coloredData := s.colorData(stencil, data)
	result, err := s.interpolate(stencil.Template, coloredData)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (s *Stenciller) findStencil(id string) (*stencil, error) {
	for _, stencil := range s.stencils {
		if stencil.ID == id {
			return stencil, nil
		}
	}
	return nil, fmt.Errorf("Unable to find stencil with id of %v", id)
}

func (s *Stenciller) colorData(stencil *stencil, data map[string]string) map[string]string {
	colored := make(map[string]string, len(data))
	for key, val := range data {
		if col, ok := stencil.Colors[key]; ok {
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
	log.Infof("Unable to set [value=%v] with color [%v]", val, col)
	return val
}

func (s *Stenciller) interpolate(tmpl string, data map[string]string) (string, error) {
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
