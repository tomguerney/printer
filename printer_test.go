package printer

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockWriter is a mock io.Writer for testing
type MockWriter struct {
	mock.Mock
}

func (m *MockWriter) Write(p []byte) (n int, err error) {
	args := m.Called(string(p))
	return args.Int(0), args.Error(1)
}

type PrinterSuite struct {
	suite.Suite
	Writer *MockWriter
}

func (suite *PrinterSuite) SetupTest() {
	suite.Writer = new(MockWriter)
	suite.Writer.On("Write", mock.Anything).Return(0, nil)
	setter.Writer = suite.Writer
}

func (suite *PrinterSuite) TestMessage() {
	Msg("Hello")
	expected := "Hello\n"
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

func (suite *PrinterSuite) TestSMessage() {
	msg := "Hello"
	expected := "Hello\n"
	actual := SMsg(msg)
	suite.Equal(expected, actual)
}

func (suite *PrinterSuite) TestMessageWithArgument() {
	msg := "Hello %v"
	args := "World"
	expected := "Hello World\n"
	Msg(msg, args)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

func (suite *PrinterSuite) TestSMessageWithArgument() {
	msg := "Hello %v"
	args := "World"
	expected := "Hello World\n"
	actual := SMsg(msg, args)
	suite.Equal(expected, actual)
}

func (suite *PrinterSuite) TestError() {
	errMsg := "error message"
	expected := "Error: error message\n"
	Error(errMsg)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

func (suite *PrinterSuite) TestErrorWithArgument() {
	errMsg := "%v message"
	args := "error"
	expected := "Error: error message\n"
	Error(errMsg, args)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

func (suite *PrinterSuite) TestSError() {
	errMsg := "error message"
	expected := "Error: error message\n"
	actual := SError(errMsg)
	suite.Equal(expected, actual)
}

func (suite *PrinterSuite) TestSErrorWithArgument() {
	errMsg := "%v message"
	args := "error"
	expected := "Error: error message\n"
	actual := SError(errMsg, args)
	suite.Equal(expected, actual)
}

func (suite *PrinterSuite) TestTabulate() {
	table := [][]string{
		{"The", "first", "row"},
		{"This", "is", "another", "row"},
		{"The", "tertiary", "row"},
	}

	expected := []string{
		"The     first       row",
		"This    is          another    row",
		"The     tertiary    row",
	}
	Tabulate(table)
	suite.Writer.AssertCalled(suite.T(), "Write", expected[0])
	suite.Writer.AssertCalled(suite.T(), "Write", expected[1])
	suite.Writer.AssertCalled(suite.T(), "Write", expected[2])
	suite.Writer.AssertNumberOfCalls(suite.T(), "Write", 3)
}
func (suite *PrinterSuite) TestSTabulate() {
	table := [][]string{
		{"The", "first", "row"},
		{"This", "is", "another", "row"},
		{"The", "tertiary", "row"},
	}

	expected := []string{
		"The     first       row",
		"This    is          another    row",
		"The     tertiary    row",
	}
	actual := STabulate(table)
	suite.Equal(expected, actual)
}

func (suite *PrinterSuite) TestAddStencilWithExistingID() {
	id := "test-id"
	template := "{{ .test }} template"
	colors := map[string]string{
		"test": "red",
	}
	err := AddStencil(id, template, colors)
	suite.NoError(err)
	err = AddStencil(id, template, colors)
	suite.Errorf(err, "Stencil with ID %v already exists")
}

func (suite *PrinterSuite) TestUseStencil() {
	id := "test-id"
	template := "{{ .test }} template"
	colors := map[string]string{
		"test": "red",
	}
	AddStencil(id, template, colors)
	data := map[string]string{
		"test": "value",
	}
	expected := "value template"
	err := UseStencil(id, data)
	suite.NoError(err)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

// func (suite *PrinterSuite) TestFindStencil() {
// 	stencil1 := &stencil{ID: "1"}
// 	stencil2 := &stencil{ID: "2"}
// 	suite.Stenciller.stencils = []*stencil{stencil1, stencil2}
// 	actual, err := suite.Stenciller.findStencil("1")
// 	suite.NoError(err)
// 	suite.Equal(stencil1, actual)
// 	suite.NotEqual(stencil2, actual)
// }

// func (suite *PrinterSuite) TestNotFindStencil() {
// 	stencil1 := &stencil{ID: "1"}
// 	stencil2 := &stencil{ID: "2"}
// 	suite.Stenciller.stencils = []*stencil{stencil1, stencil2}
// 	actual, err := suite.Stenciller.findStencil("3")
// 	suite.Errorf(err, "Unable to find stencil with id of 3")
// 	suite.Nil(actual)
// }

// func (suite *PrinterSuite) TestColorData() {
// 	expected := map[string]string{
// 		"key1": "blueValue",
// 		"key2": "greenValue",
// 	}
// 	stencil := &stencil{
// 		ID: "1",
// 		Colors: map[string]string{
// 			"key1": "blue",
// 			"key2": "green",
// 		},
// 	}
// 	data := map[string]string{
// 		"key1": "value1",
// 		"key2": "value2",
// 	}
// 	suite.Colorer.On("Color", "value1", "blue").Return("blueValue", true)
// 	suite.Colorer.On("Color", "value2", "green").Return("greenValue", true)
// 	actual := suite.Stenciller.colorData(stencil, data)
// 	suite.Equal(expected, actual)
// }

// func (suite *PrinterSuite) TestColorDataWithValueWithoutColorDefinition() {
// 	expected := map[string]string{
// 		"key1": "blueValue1",
// 		"key2": "value2",
// 		"key3": "greenValue3",
// 	}
// 	stencil := &stencil{
// 		ID: "1",
// 		Colors: map[string]string{
// 			"key1": "blue",
// 			"key3": "green",
// 		},
// 	}
// 	data := map[string]string{
// 		"key1": "value1",
// 		"key2": "value2",
// 		"key3": "value3",
// 	}
// 	suite.Colorer.On("Color", "value1", "blue").Return("blueValue1", true)
// 	suite.Colorer.On("Color", "value3", "green").Return("greenValue3", true)
// 	actual := suite.Stenciller.colorData(stencil, data)
// 	suite.Equal(expected, actual)
// }

// func (suite *PrinterSuite) TestColorDataWithNonExistantColor() {
// 	expected := map[string]string{
// 		"key1": "blueValue1",
// 		"key2": "value2",
// 		"key3": "value3",
// 	}
// 	stencil := &stencil{
// 		ID: "1",
// 		Colors: map[string]string{
// 			"key1": "blue",
// 			"key3": "notacolor",
// 		},
// 	}
// 	data := map[string]string{
// 		"key1": "value1",
// 		"key2": "value2",
// 		"key3": "value3",
// 	}
// 	suite.Colorer.On("Color", "value1", "blue").Return("blueValue1", true)
// 	suite.Colorer.On("Color", "value3", "notacolor").Return("", false)
// 	actual := suite.Stenciller.colorData(stencil, data)
// 	suite.Equal(expected, actual)
// }

// func (suite *PrinterSuite) TestInterpolate() {
// 	data := map[string]string{
// 		"Field1": "value1",
// 		"field2": "value2",
// 	}
// 	tmpl := "abc {{ .field2 }} def {{.Field1}}"
// 	expected := "abc value2 def value1"
// 	actual, err := suite.Stenciller.interpolate(tmpl, data)
// 	suite.NoError(err)
// 	suite.Equal(expected, actual)
// }

// func (suite *PrinterSuite) TestInterpolateWithNonExistantKey() {
// 	data := map[string]string{
// 		"Field1": "value1",
// 		"field2": "value2",
// 	}
// 	tmpl := "abc {{ .field2 }} def {{.Field3}}"
// 	expected := "abc value2 def "
// 	actual, err := suite.Stenciller.interpolate(tmpl, data)
// 	suite.NoError(err)
// 	suite.Equal(expected, actual)
// }

// func (suite *PrinterSuite) TestInterpolateWithExtraMapKey() {
// 	data := map[string]string{
// 		"Field1": "value1",
// 		"field2": "value2",
// 		"Field3": "value3",
// 	}
// 	tmpl := "abc {{ .field2 }} def {{.Field1}}"
// 	expected := "abc value2 def value1"
// 	actual, err := suite.Stenciller.interpolate(tmpl, data)
// 	suite.NoError(err)
// 	suite.Equal(expected, actual)
// }

func TestPrinterSuite(t *testing.T) {
	suite.Run(t, new(PrinterSuite))
}
