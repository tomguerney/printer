package stenciller

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type StencillerSuite struct {
	suite.Suite
	Stenciller *Stenciller
}

func (suite *StencillerSuite) SetupTest() {
	suite.Stenciller = &Stenciller{}
}

func (suite *StencillerSuite) TestAddStencil() {
	suite.Empty(suite.Stenciller.stencils)
	id := "test-id"
	template := "{{ .Test }} template"
	colors := map[string]string{
		"test": "red",
	}
	suite.Stenciller.AddStencil(id, template, colors)
	suite.Len(suite.Stenciller.stencils, 1)
}



func TestStencillerSuite(t *testing.T) {
	suite.Run(t, new(StencillerSuite))
}
