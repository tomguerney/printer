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
	suite.Empty(suite.Stenciller.templateStencils)
	id := "test-id"
	template := "{{ .test }} template"
	colors := map[string]string{
		"test": "red",
	}
	err := suite.Stenciller.AddTemplateStencil(id, template, colors)
	suite.NoError(err)
	suite.Len(suite.Stenciller.templateStencils, 1)
}

func (suite *StencillerSuite) TestAddTmplStencilWithExistingID() {
	stencil := &TemplateStencil{ID: "test-id",
		Template: "{{ .test }} template",
		Colors: map[string]string{
			"test": "red",
		}}
	suite.Stenciller.templateStencils =
		append(suite.Stenciller.templateStencils, stencil)
	suite.Len(suite.Stenciller.templateStencils, 1)
	err := suite.Stenciller.AddTemplateStencil(stencil.ID, "{{ .Template }}", nil)
	suite.Errorf(err, "Template Stencil with ID test-id already exists")
}

func (suite *StencillerSuite) TestAddTmplStencilWithEmptyID() {
	err := suite.Stenciller.AddTemplateStencil("", "{{ .Template }}", nil)
	suite.Errorf(err, "Stencil ID may not be empty")
}

func (suite *StencillerSuite) TestAddTableStencil() {
	suite.Empty(suite.Stenciller.tableStencils)
	id := "test-id"
	colors := map[string]string{
		"test": "red",
	}
	headers := []string{"header1", "header2"}
	columnOrder := []string{"key1", "key2"}
	err := suite.Stenciller.AddTableStencil(id, headers, columnOrder, colors)
	suite.NoError(err)
	suite.Len(suite.Stenciller.tableStencils, 1)
}

func (suite *StencillerSuite) TestAddTableStencilWithExistingID() {
	stencil := &TableStencil{ID: "test-id",
		Headers: []string{"header1", "header2"},
		Colors: map[string]string{
			"test": "red",
		}}
	suite.Stenciller.tableStencils =
		append(suite.Stenciller.tableStencils, stencil)
	suite.Len(suite.Stenciller.tableStencils, 1)
	err := suite.Stenciller.AddTableStencil(stencil.ID, nil, nil, nil)
	suite.Errorf(err, "Table Stencil with ID test-id already exists")
}

func (suite *StencillerSuite) TestAddTableStencilWithEmptyID() {
	err := suite.Stenciller.AddTableStencil("", nil, nil, nil)
	suite.Errorf(err, "Stencil ID may not be empty")
}

func (suite *StencillerSuite) TestTmplStencil() {
	id := "test-id"
	template := "{{ .test }} template"
	colors := map[string]string{
		"test": "red",
	}
	suite.Stenciller.AddTemplateStencil(id, template, colors)
	data := map[string]string{
		"test": "value",
	}
	suite.Colorer.On("Color", "value", "red").Return("redValue", true)
	actual, err := suite.Stenciller.UseTemplateStencil(id, data)
	suite.NoError(err)
	expected := "redValue template"
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestTableStencil() {
	id := "test-id"
	colors := map[string]string{
		"key2": "red",
	}
	columnOrder := []string{"key1", "key2"}
	suite.Stenciller.AddTableStencil(id, nil, columnOrder, colors)
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
	actual, err := suite.Stenciller.UseTableStencil(id, data)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestTableStencilWithHeaders() {
	id := "test-id"
	colors := map[string]string{
		"key2": "red",
	}
	headers := []string{"header1", "header2"}
	columnOrder := []string{"key1", "key2"}
	suite.Stenciller.AddTableStencil(id, headers, columnOrder, colors)
	data := []map[string]string{{
		"key1": "value1a",
		"key2": "value2a",
	}, {
		"key1": "value1b",
		"key2": "value2b",
	}}
	expected := [][]string{
		{"header1", "header2"},
		{"-------", "-------"},
		{"value1a", "redValue"},
		{"value1b", "redValue"},
	}
	suite.Colorer.On("Color", mock.Anything, "red").Return("redValue", true)
	actual, err := suite.Stenciller.UseTableStencil(id, data)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestTableStencilWithMoreHeadersThanDataCols() {
	id := "test-id"
	colors := map[string]string{
		"key2": "red",
	}
	headers := []string{"header1", "header2", "header3"}
	columnOrder := []string{"key1", "key2"}
	suite.Stenciller.AddTableStencil(id, headers, columnOrder, colors)
	data := []map[string]string{{
		"key1": "value1a",
		"key2": "value2a",
	}, {
		"key1": "value1b",
		"key2": "value2b",
	}}
	expected := [][]string{
		{"header1", "header2", "header3"},
		{"-------", "-------", "-------"},
		{"value1a", "redValue"},
		{"value1b", "redValue"},
	}
	suite.Colorer.On("Color", mock.Anything, "red").Return("redValue", true)
	actual, err := suite.Stenciller.UseTableStencil(id, data)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestTableStencilWithHeadersLongerThanItems() {
	id := "test-id"
	colors := map[string]string{
		"key2": "red",
	}
	headers := []string{"header1", "header2IsQuiteLong"}
	columnOrder := []string{"key1", "key2"}
	suite.Stenciller.AddTableStencil(id, headers, columnOrder, colors)
	data := []map[string]string{{
		"key1": "value1a",
		"key2": "value2a",
	}, {
		"key1": "value1b",
		"key2": "value2b",
	}}
	expected := [][]string{
		{"header1", "header2IsQuiteLong"},
		{"-------", "------------------"},
		{"value1a", "redValue"},
		{"value1b", "redValue"},
	}
	suite.Colorer.On("Color", mock.Anything, "red").Return("redValue", true)
	actual, err := suite.Stenciller.UseTableStencil(id, data)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestTableStencilWithHeadersAndOneRow() {
	id := "test-id"
	colors := map[string]string{
		"key2": "red",
	}
	headers := []string{"header1", "header2"}
	columnOrder := []string{"key1", "key2"}
	suite.Stenciller.AddTableStencil(id, headers, columnOrder, colors)
	data := []map[string]string{{
		"key1": "value1a",
		"key2": "value2a",
	}}
	expected := [][]string{
		{"header1", "header2"},
		{"-------", "-------"},
		{"value1a", "redValue"},
	}
	suite.Colorer.On("Color", mock.Anything, "red").Return("redValue", true)
	actual, err := suite.Stenciller.UseTableStencil(id, data)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestTableStencilWithOutOfOrderData() {
	id := "test-id"
	colors := map[string]string{
		"key2": "red",
	}
	columnOrder := []string{"key1", "key2"}
	suite.Stenciller.AddTableStencil(id, nil, columnOrder, colors)
	data := []map[string]string{{
		"key1": "value1a",
		"key2": "value2a",
	}, {
		"key2": "value2b",
		"key1": "value1b",
	}}
	expected := [][]string{
		{"value1a", "redValue"},
		{"value1b", "redValue"},
	}
	suite.Colorer.On("Color", mock.Anything, "red").Return("redValue", true)
	actual, err := suite.Stenciller.UseTableStencil(id, data)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestTableStencilWithDataNotInColumnOrder() {
	id := "test-id"
	colors := map[string]string{
		"key2": "red",
	}
	columnOrder := []string{"key1", "key2"}
	suite.Stenciller.AddTableStencil(id, nil, columnOrder, colors)
	data := []map[string]string{{
		"key1": "value1a",
		"key3": "this should not appear",
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
	actual, err := suite.Stenciller.UseTableStencil(id, data)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestTableStencilWithColumnOrderNotInData() {
	id := "test-id"
	colors := map[string]string{
		"key2": "red",
	}
	columnOrder := []string{"key1", "notinsomedata", "key2"}
	suite.Stenciller.AddTableStencil(id, nil, columnOrder, colors)
	data := []map[string]string{{
		"key1": "value1a",
		"key2": "value2a",
	}, {
		"key1":          "value1b",
		"notinsomedata": "anothervalue",
		"key2":          "value2b",
	}}
	expected := [][]string{
		{"value1a", "", "redValue"},
		{"value1b", "anothervalue", "redValue"},
	}
	suite.Colorer.On("Color", mock.Anything, "red").Return("redValue", true)
	actual, err := suite.Stenciller.UseTableStencil(id, data)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestFindTmplStencil() {
	stencil1 := &TemplateStencil{ID: "1"}
	stencil2 := &TemplateStencil{ID: "2"}
	suite.Stenciller.templateStencils = []*TemplateStencil{stencil1, stencil2}
	actual, err := suite.Stenciller.findTemplateStencil("1")
	suite.NoError(err)
	suite.Equal(stencil1, actual)
	suite.NotEqual(stencil2, actual)
}

func (suite *StencillerSuite) TestNotFindTmplStencil() {
	stencil1 := &TemplateStencil{ID: "1"}
	stencil2 := &TemplateStencil{ID: "2"}
	suite.Stenciller.templateStencils = []*TemplateStencil{stencil1, stencil2}
	actual, err := suite.Stenciller.findTemplateStencil("3")
	suite.Errorf(err, "Unable to find stencil with id of 3")
	suite.Nil(actual)
}

func (suite *StencillerSuite) TestColorMap() {
	expected := map[string]string{
		"key1": "blueValue",
		"key2": "greenValue",
	}
	stencil := &TemplateStencil{
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
	actual := suite.Stenciller.colorMap(stencil.Colors, data)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestColorDataWithValueWithoutColorDefinition() {
	expected := map[string]string{
		"key1": "blueValue1",
		"key2": "value2",
		"key3": "greenValue3",
	}
	stencil := &TemplateStencil{
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
	actual := suite.Stenciller.colorMap(stencil.Colors, data)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestColorDataWithNonExistantColor() {
	expected := map[string]string{
		"key1": "blueValue1",
		"key2": "value2",
		"key3": "value3",
	}
	stencil := &TemplateStencil{
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
	actual := suite.Stenciller.colorMap(stencil.Colors, data)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestInterpolate() {
	data := map[string]string{
		"Field1": "value1",
		"field2": "value2",
	}
	tmpl := "abc {{ .field2 }} def {{.Field1}}"
	expected := "abc value2 def value1"
	actual, err := interpolate(tmpl, data)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestInterpolateWithNonExistantKey() {
	data := map[string]string{
		"Field1": "value1",
		"field2": "value2",
	}
	tmpl := "abc {{ .field2 }} def {{.Field3}}"
	expected := "abc value2 def <no value>"
	actual, err := interpolate(tmpl, data)
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
	actual, err := interpolate(tmpl, data)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func TestStencillerSuite(t *testing.T) {
	suite.Run(t, new(StencillerSuite))
}
