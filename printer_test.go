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
	s.OutWriter = suite.Writer
}

func (suite *PrinterSuite) TestMessage() {
	Out("Hello")
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

func (suite *PrinterSuite) TestTabulate() {
	table := [][]string{
		{"The", "first", "row"},
		{"This", "is", "another", "row"},
		{"The", "tertiary", "row"},
	}

	expected := []string{
		"The     first       row\n",
		"This    is          another    row\n",
		"The     tertiary    row\n",
	}
	Tabulate(table)
	suite.Writer.AssertCalled(suite.T(), "Write", expected[0])
	suite.Writer.AssertCalled(suite.T(), "Write", expected[1])
	suite.Writer.AssertCalled(suite.T(), "Write", expected[2])
	suite.Writer.AssertNumberOfCalls(suite.T(), "Write", 3)
}
func (suite *PrinterSuite) TestAddStencilWithExistingID() {
	id := "test-id"
	template := "{{ .test }} template"
	colors := map[string]string{
		"test": "red",
	}
	err := AddTmplStencil(id, template, colors)
	suite.NoError(err)
	err = AddTmplStencil(id, template, colors)
	suite.Errorf(err, "Stencil with ID %v already exists")
}

func (suite *PrinterSuite) TestTmplStencil() {
	id := "test-id"
	template := "{{ .test }} template"
	colors := map[string]string{
		"test": "red",
	}
	AddTmplStencil(id, template, colors)
	data := map[string]string{
		"test": "value",
	}
	expected := "value template\n"
	err := TmplStencil(id, data)
	suite.NoError(err)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

func (suite *PrinterSuite) TestTableStencil() {
	id := "test-id"
	headers := []string{"Header1", "This is header 2", "h3"}
	columnOrder := []string{"Key1", "Key3", "Key2"}
	colors := map[string]string{"Key3": "green"}
	AddTableStencil(id, headers, columnOrder, colors)
	data := []map[string]string{
		{
			"Key1": "value1a",
			"Key2": "This is value 2a",
			"Key3": "v3a",
		},
		{
			"Key1": "value1b",
			"Key2": "value 2b",
			"Key3": "v3b",
		},
	}
	expected := []string{
		"Header1    This is header 2    h3\n",
		"-------    ----------------    ----------------\n",
		"value1a    v3a                 This is value 2a\n",
		"value1b    v3b                 value 2b\n",
	}
	TableStencil(id, data)
	suite.Writer.AssertCalled(suite.T(), "Write", expected[0])
	suite.Writer.AssertCalled(suite.T(), "Write", expected[1])
	suite.Writer.AssertCalled(suite.T(), "Write", expected[2])
	suite.Writer.AssertCalled(suite.T(), "Write", expected[3])
	suite.Writer.AssertNumberOfCalls(suite.T(), "Write", 4)
}

func TestPrinterSuite(t *testing.T) {
	suite.Run(t, new(PrinterSuite))
}
