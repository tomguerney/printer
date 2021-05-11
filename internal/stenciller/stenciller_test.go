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

func (m *MockColorer) Color(text, color string) (string, bool) {
	args := m.Called(text, color)
	return args.String(0), args.Bool(1)
}

func (suite *StencillerSuite) SetupTest() {
	suite.Colorer = new(MockColorer)
	suite.Stenciller = &Stenciller{colorer: suite.Colorer}
}

func (suite *StencillerSuite) TestAddStencil() {
	suite.Empty(suite.Stenciller.stencils)
	id := "test-id"
	template := "{{ .test }} template"
	colors := map[string]string{
		"test": "red",
	}
	err := suite.Stenciller.AddStencil(id, template, colors)
	suite.NoError(err)
	suite.Len(suite.Stenciller.stencils, 1)
}

func (suite *StencillerSuite) TestAddStencilWithExistingID() {
	stencil := &stencil{ID: "test-id",
		Template: "{{ .test }} template",
		Colors: map[string]string{
			"test": "red",
		}}
	suite.Stenciller.stencils = append(suite.Stenciller.stencils, stencil)
	suite.Len(suite.Stenciller.stencils, 1)
	err := suite.Stenciller.AddStencil(stencil.ID, "{{ .Template }}", nil)
	suite.Errorf(err, "Stencil with ID test-id already exists")
}

func (suite *StencillerSuite) TestUseStencil() {
	id := "test-id"
	template := "{{ .test }} template"
	colors := map[string]string{
		"test": "red",
	}
	suite.Stenciller.AddStencil(id, template, colors)
	data := map[string]string {
		"test": "value",
	}
	suite.Colorer.On("Color", "value", "red").Return("redValue", true)
	actual, err := suite.Stenciller.UseStencil(id, data)
	suite.NoError(err)
	expected := "redValue template\n"
	suite.Equal(expected, actual)
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
	suite.Colorer.On("Color", "value1", "blue").Return("blueValue", true)
	suite.Colorer.On("Color", "value2", "green").Return("greenValue", true)
	actual := suite.Stenciller.colorData(stencil, data)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestColorDataWithValueWithoutColorDefinition() {
	expected := map[string]string{
		"key1": "blueValue1",
		"key2": "value2",
		"key3": "greenValue3",
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
		"key2": "value2",
		"key3": "value3",
	}
	suite.Colorer.On("Color", "value1", "blue").Return("blueValue1", true)
	suite.Colorer.On("Color", "value3", "green").Return("greenValue3", true)
	actual := suite.Stenciller.colorData(stencil, data)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestColorDataWithNonExistantColor() {
	expected := map[string]string{
		"key1": "blueValue1",
		"key2": "value2",
		"key3": "value3",
	}
	stencil := &stencil{
		ID: "1",
		Colors: map[string]string{
			"key1": "blue",
			"key3": "notacolor",
		},
	}
	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	suite.Colorer.On("Color", "value1", "blue").Return("blueValue1", true)
	suite.Colorer.On("Color", "value3", "notacolor").Return("", false)
	actual := suite.Stenciller.colorData(stencil, data)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestInterpolate() {
	data := map[string]string{
		"Field1": "value1",
		"field2": "value2",
	}
	tmpl := "abc {{ .field2 }} def {{.Field1}}"
	expected := "abc value2 def value1"
	actual, err := suite.Stenciller.interpolate(tmpl, data)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestInterpolateWithNonExistantKey() {
	data := map[string]string{
		"Field1": "value1",
		"field2": "value2",
	}
	tmpl := "abc {{ .field2 }} def {{.Field3}}"
	expected := "abc value2 def "
	actual, err := suite.Stenciller.interpolate(tmpl, data)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestInterpolateWithExtraMapKey() {
	data := map[string]string{
		"Field1": "value1",
		"field2": "value2",
		"Field3": "value3",
	}
	tmpl := "abc {{ .field2 }} def {{.Field1}}"
	expected := "abc value2 def value1"
	actual, err := suite.Stenciller.interpolate(tmpl, data)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func TestStencillerSuite(t *testing.T) {
	suite.Run(t, new(StencillerSuite))
}
