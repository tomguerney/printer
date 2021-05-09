package stenciller

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type StencillerSuite struct {
	suite.Suite
	Stenciller *Stenciller
	Colorer    *MockColorer
}

type MockColorer struct {
	mock.Mock
}

func (m *MockColorer) Color(text, color string) (string, error) {
	args := m.Called(text, color)
	return args.String(0), args.Error(1)
}

func (suite *StencillerSuite) SetupTest() {
	suite.Colorer = new(MockColorer)
	suite.Stenciller = &Stenciller{colorer: suite.Colorer}
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
func (suite *StencillerSuite) TestFindStencil() {
	stencil1 := &stencil{ID: "1"}
	stencil2 := &stencil{ID: "2"}
	suite.Stenciller.stencils = []*stencil{stencil1, stencil2}
	actual, err := suite.Stenciller.findStencil("1")
	suite.NoError(err)
	suite.Equal(stencil1, actual)
	suite.NotEqual(stencil2, actual)
}

func (suite *StencillerSuite) TestNotFindStencil() {
	stencil1 := &stencil{ID: "1"}
	stencil2 := &stencil{ID: "2"}
	suite.Stenciller.stencils = []*stencil{stencil1, stencil2}
	actual, err := suite.Stenciller.findStencil("3")
	suite.Errorf(err, "Unable to find stencil with id of 3")
	suite.Nil(actual)
}

func (suite *StencillerSuite) TestColorData() {
	expected := map[string]string{
		"key1": "blueValue",
		"key2": "greenValue",
	}
	stencil := &stencil{
		ID: "1",
		Colors: map[string]string{
			"key1": "blue",
			"key2": "green",
		},
	}
	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	suite.Colorer.On("Color", "value1", "blue").Return("blueValue", nil)
	suite.Colorer.On("Color", "value2", "green").Return("greenValue", nil)
	actual := suite.Stenciller.colorData(stencil, data)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestColorDataWithValueWithoutColorDefinition() {
	expected := map[string]string{
		"key1": "blueValue",
		"key2": "originalValue",
		"key3": "greenValue",
	}
	stencil := &stencil{
		ID: "1",
		Colors: map[string]string{
			"key1": "blue",
			"key3": "green",
		},
	}
	data := map[string]string{
		"key1": "value1",
		"key2": "originalValue",
		"key3": "value3",
	}
	suite.Colorer.On("Color", "value1", "blue").Return("blueValue", nil)
	suite.Colorer.On("Color", "value3", "green").Return("greenValue", nil)
	actual := suite.Stenciller.colorData(stencil, data)
	suite.Equal(expected, actual)
}

func TestStencillerSuite(t *testing.T) {
	suite.Run(t, new(StencillerSuite))
}
