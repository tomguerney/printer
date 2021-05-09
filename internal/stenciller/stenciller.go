package stenciller

import (
	"fmt"

	"github.com/apex/log"
	"github.com/tomguerney/printer/internal/domain"
)

// Stenciller takes structs and formats their field values to predefined templates and colors
type Stenciller struct {
	stencils []*stencil
	colorer  domain.Colorer
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
	//TODO: validate input: id is not duplicated, colors exist, template is valid
	s.stencils = append(s.stencils, &stencil{id, template, colors})
	return nil
}

// UseStencil applies an interface to the stencil with the passed ID
func (s *Stenciller) UseStencil(id string, data map[string]string) (string, error) {
	stencil, err := s.findStencil(id)
	if err != nil {
		return "", err
	}
	_ = s.colorData(stencil, data)
	return "", nil
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
			coloredVal, err := s.colorer.Color(val, col)
			if err != nil {
				log.Infof("Unable to set [key=%v] [value=%v] with color [%v]")
				colored[key] = val
			} else {
				colored[key] = coloredVal
			}
		} else {
			colored[key] = val
		}
	}
	return colored
}
