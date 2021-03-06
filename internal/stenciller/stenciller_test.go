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
	stencil := &TemplateStencil{
		ID:       id,
		Template: template,
		Colors:   colors,
	}
	err := suite.Stenciller.AddTemplateStencil(stencil)
	suite.NoError(err)
	suite.Len(suite.Stenciller.templateStencils, 1)
}

func (suite *StencillerSuite) TestAddTmplStencilWithExistingID() {
	stencil1 := &TemplateStencil{
		ID:       "test-id",
		Template: "{{ .test }} template",
		Colors: map[string]string{
			"test": "red",
		}}
	suite.Stenciller.templateStencils = append(suite.Stenciller.templateStencils, stencil1)
	stencil2 := &TemplateStencil{
		ID:       stencil1.ID,
		Template: "{{ .Template }}",
	}
	suite.Len(suite.Stenciller.templateStencils, 1)
	err := suite.Stenciller.AddTemplateStencil(stencil2)
	suite.Errorf(err, "Template Stencil with ID test-id already exists")
}

func (suite *StencillerSuite) TestAddTmplStencilWithEmptyID() {
	stencil := &TemplateStencil{
		Template: "{{ .Template }}",
	}
	err := suite.Stenciller.AddTemplateStencil(stencil)
	suite.Errorf(err, "Stencil ID may not be empty")
}

func (suite *StencillerSuite) TestAddTableStencil() {
	suite.Empty(suite.Stenciller.tableStencils)
	stencil := &TableStencil{
		ID:          "test-id",
		Headers:     []string{"header1", "header2"},
		Colors:      map[string]string{"test": "red"},
		ColumnOrder: []string{"key1", "key2"},
	}
	err := suite.Stenciller.AddTableStencil(stencil)
	suite.NoError(err)
	suite.Len(suite.Stenciller.tableStencils, 1)
}

func (suite *StencillerSuite) TestAddTableStencilWithExistingID() {
	stencil := &TableStencil{
		ID:      "test-id",
		Headers: []string{"header1", "header2"},
		Colors:  map[string]string{"test": "red"},
	}
	suite.Stenciller.tableStencils =
		append(suite.Stenciller.tableStencils, stencil)
	suite.Len(suite.Stenciller.tableStencils, 1)
	err := suite.Stenciller.AddTableStencil(&TableStencil{ID: stencil.ID})
	suite.Errorf(err, "Table Stencil with ID test-id already exists")
}

func (suite *StencillerSuite) TestAddTableStencilWithEmptyID() {
	stencil := &TableStencil{ID: ""}
	err := suite.Stenciller.AddTableStencil(stencil)
	suite.Errorf(err, "Stencil ID may not be empty")
}

func (suite *StencillerSuite) TestTmplStencil() {
	stencil := &TemplateStencil{
		ID:       "test-id",
		Template: "{{ .test }} template",
		Colors: map[string]string{
			"test": "red",
		},
	}
	suite.Stenciller.AddTemplateStencil(stencil)
	data := map[string]string{
		"test": "value",
	}
	suite.Colorer.On("Color", "value", "red").Return("redValue", true)
	actual, err := suite.Stenciller.UseTemplateStencil(stencil.ID, data)
	suite.NoError(err)
	expected := "redValue template"
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestTableStencil() {
	stencil := &TableStencil{
		ID: "test-id",
		Colors: map[string]string{
			"key2": "red",
		},
		ColumnOrder: []string{"key1", "key2"},
	}
	suite.Stenciller.AddTableStencil(stencil)
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
	actual, err := suite.Stenciller.UseTableStencil(stencil.ID, data)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestTableStencilWithHeaders() {
	stencil := &TableStencil{
		ID: "test-id",
		Colors: map[string]string{
			"key2": "red",
		},
		ColumnOrder: []string{"key1", "key2"},
		Headers:     []string{"header1", "header2"},
	}
	suite.Stenciller.AddTableStencil(stencil)
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
	actual, err := suite.Stenciller.UseTableStencil(stencil.ID, data)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestTableStencilWithMoreHeadersThanDataCols() {
	stencil := &TableStencil{
		ID: "test-id",
		Colors: map[string]string{
			"key2": "red",
		},
		ColumnOrder: []string{"key1", "key2"},
		Headers:     []string{"header1", "header2", "header3"},
	}
	suite.Stenciller.AddTableStencil(stencil)
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
	actual, err := suite.Stenciller.UseTableStencil(stencil.ID, data)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestTableStencilWithHeadersLongerThanItems() {
	stencil := &TableStencil{
		ID: "test-id",
		Colors: map[string]string{
			"key2": "red",
		},
		ColumnOrder: []string{"key1", "key2"},
		Headers:     []string{"header1", "header2IsQuiteLong"},
	}
	suite.Stenciller.AddTableStencil(stencil)
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
	actual, err := suite.Stenciller.UseTableStencil(stencil.ID, data)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestTableStencilWithHeadersAndOneRow() {
	stencil := &TableStencil{
		ID: "test-id",
		Colors: map[string]string{
			"key2": "red",
		},
		ColumnOrder: []string{"key1", "key2"},
		Headers:     []string{"header1", "header2"},
	}
	suite.Stenciller.AddTableStencil(stencil)
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
	actual, err := suite.Stenciller.UseTableStencil(stencil.ID, data)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestTableStencilWithOutOfOrderData() {
	stencil := &TableStencil{
		ID:          "test-id",
		Colors:      map[string]string{"key2": "red"},
		ColumnOrder: []string{"key1", "key2"},
	}
	suite.Stenciller.AddTableStencil(stencil)
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
	actual, err := suite.Stenciller.UseTableStencil(stencil.ID, data)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestTableStencilWithDataNotInColumnOrder() {
	stencil := &TableStencil{
		ID:          "test-id",
		Colors:      map[string]string{"key2": "red"},
		ColumnOrder: []string{"key1", "key2"},
	}
	suite.Stenciller.AddTableStencil(stencil)
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
	actual, err := suite.Stenciller.UseTableStencil(stencil.ID, data)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *StencillerSuite) TestTableStencilWithColumnOrderNotInData() {
	stencil := &TableStencil{
		ID:          "test-id",
		Colors:      map[string]string{"key2": "red"},
		ColumnOrder: []string{"key1", "notinsomedata", "key2"},
	}
	suite.Stenciller.AddTableStencil(stencil)
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
	actual, err := suite.Stenciller.UseTableStencil(stencil.ID, data)
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

func TestStencillerSuite(t *testing.T) {
	suite.Run(t, new(StencillerSuite))
}
