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

func (suite *StencillerSuite) TestAddTmplStencil() {
	suite.Empty(suite.Stenciller.tmplStencils)
	id := "test-id"
	template := "{{ .test }} template"
	colors := map[string]string{
		"test": "red",
	}
	err := suite.Stenciller.AddTmplStencil(id, template, colors)
	suite.NoError(err)
	suite.Len(suite.Stenciller.tmplStencils, 1)
}

func (suite *StencillerSuite) TestAddTmplStencilWithExistingID() {
	stencil := &tmplStencil{ID: "test-id",
		Template: "{{ .test }} template",
		Colors: map[string]string{
			"test": "red",
		}}
	suite.Stenciller.tmplStencils =
		append(suite.Stenciller.tmplStencils, stencil)
	suite.Len(suite.Stenciller.tmplStencils, 1)
	err := suite.Stenciller.AddTmplStencil(stencil.ID, "{{ .Template }}", nil)
	suite.Errorf(err, "Template Stencil with ID test-id already exists")
}

func (suite *StencillerSuite) TestAddTmplStencilWithEmptyID() {
	err := suite.Stenciller.AddTmplStencil("", "{{ .Template }}", nil)
	suite.Errorf(err, "Stencil ID may not be empty")
}

func (suite *StencillerSuite) TestAddTableStencil() {
	suite.Empty(suite.Stenciller.tableStencils)
	id := "test-id"
	colors := map[string]string{
		"test": "red",
	}
	headers := []string{"header1", "header2"}
	err := suite.Stenciller.AddTableStencil(id, headers, colors)
	suite.NoError(err)
	suite.Len(suite.Stenciller.tableStencils, 1)
}

func (suite *StencillerSuite) TestAddTableStencilWithExistingID() {
	stencil := &tableStencil{ID: "test-id",
		Headers: []string{"header1", "header2"},
		Colors: map[string]string{
			"test": "red",
		}}
	suite.Stenciller.tableStencils =
		append(suite.Stenciller.tableStencils, stencil)
	suite.Len(suite.Stenciller.tableStencils, 1)
	err := suite.Stenciller.AddTableStencil(stencil.ID, nil, nil)
	suite.Errorf(err, "Table Stencil with ID test-id already exists")
}

func (suite *StencillerSuite) TestAddTableStencilWithEmptyID() {
	err := suite.Stenciller.AddTableStencil("", nil, nil)
	suite.Errorf(err, "Stencil ID may not be empty")
}

func (suite *StencillerSuite) TestTmplStencil() {
	id := "test-id"
	template := "{{ .test }} template"
	colors := map[string]string{
		"test": "red",
	}
	suite.Stenciller.AddTmplStencil(id, template, colors)
	data := map[string]string{
		"test": "value",
	}
	suite.Colorer.On("Color", "value", "red").Return("redValue", true)
	actual, err := suite.Stenciller.TmplStencil(id, data)
	suite.NoError(err)
	expected := "redValue template\n"
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestTableStencil() {
	id := "test-id"
	colors := map[string]string{
		"key2": "red",
	}
	suite.Stenciller.AddTableStencil(id, nil, colors)
	data := []map[string]string{{
		"key1": "value1a",
		"key2": "value2a",
	}, {
		"key1": "value1b",
		"key2": "value2b",
	}}
	expected := [][]string{
		{"value1a", "redValue"},
		{"value1b", "redValue"},
	}
	suite.Colorer.On("Color", mock.Anything, "red").Return("redValue", true)
	actual, err := suite.Stenciller.TableStencil(id, data)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestTableStencilWithHeaders() {
	id := "test-id"
	colors := map[string]string{
		"key2": "red",
	}
	headers := []string{"header1", "header2"}
	suite.Stenciller.AddTableStencil(id, headers, colors)
	data := []map[string]string{{
		"key1": "value1a",
		"key2": "value2a",
	}, {
		"key1": "value1b",
		"key2": "value2b",
	}}
	expected := [][]string{
		{"header1", "header2"},
		{"-------", "--------"},
		{"value1a", "redValue"},
		{"value1b", "redValue"},
	}
	suite.Colorer.On("Color", mock.Anything, "red").Return("redValue", true)
	actual, err := suite.Stenciller.TableStencil(id, data)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestFindTmplStencil() {
	stencil1 := &tmplStencil{ID: "1"}
	stencil2 := &tmplStencil{ID: "2"}
	suite.Stenciller.tmplStencils = []*tmplStencil{stencil1, stencil2}
	actual, err := suite.Stenciller.findTmplStencil("1")
	suite.NoError(err)
	suite.Equal(stencil1, actual)
	suite.NotEqual(stencil2, actual)
}

func (suite *StencillerSuite) TestNotFindTmplStencil() {
	stencil1 := &tmplStencil{ID: "1"}
	stencil2 := &tmplStencil{ID: "2"}
	suite.Stenciller.tmplStencils = []*tmplStencil{stencil1, stencil2}
	actual, err := suite.Stenciller.findTmplStencil("3")
	suite.Errorf(err, "Unable to find stencil with id of 3")
	suite.Nil(actual)
}

func (suite *StencillerSuite) TestColorData() {
	expected := map[string]string{
		"key1": "blueValue",
		"key2": "greenValue",
	}
	stencil := &tmplStencil{
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
	actual := suite.Stenciller.colorData(stencil.Colors, data)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestColorDataWithValueWithoutColorDefinition() {
	expected := map[string]string{
		"key1": "blueValue1",
		"key2": "value2",
		"key3": "greenValue3",
	}
	stencil := &tmplStencil{
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
	actual := suite.Stenciller.colorData(stencil.Colors, data)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestColorDataWithNonExistantColor() {
	expected := map[string]string{
		"key1": "blueValue1",
		"key2": "value2",
		"key3": "value3",
	}
	stencil := &tmplStencil{
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
	actual := suite.Stenciller.colorData(stencil.Colors, data)
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
