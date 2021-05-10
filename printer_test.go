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

// MockFormatter is a mock formatter for testing
type MockFormatter struct {
	mock.Mock
}

func (m *MockFormatter) Msg(text string, a ...interface{}) string {
	args := m.Called(text, a)
	return args.String(0)
}

func (m *MockFormatter) Error(text string, a ...interface{}) string {
	args := m.Called(text, a)
	return args.String(0)
}

func (m *MockFormatter) Tabulate(table [][]string) []string {
	args := m.Called(table)
	return args.Get(0).([]string)
}

// MockTemplater is a mock templater for testing
type MockStenciller struct {
	mock.Mock
}

func (m *MockStenciller) AddStencil(id, template string, colors map[string]string) error {
	args := m.Called(id, template, colors)
	return args.Error(0)
}

func (m *MockStenciller) UseStencil(id string, data map[string]string) (string, error) {
	args := m.Called(id, data)
	return args.String(0), args.Error(1)
}

type PrinterSuite struct {
	suite.Suite
	Writer     *MockWriter
	Formatter  *MockFormatter
	Stenciller *MockStenciller
}

func (suite *PrinterSuite) SetupTest() {
	mockWriter := new(MockWriter)
	mockWriter.On("Write", mock.Anything).Return(0, nil)
	formatter = new(MockFormatter)
	stenciller = new(MockStenciller)
}

func (suite *PrinterSuite) TestMessage() {
	msg := "test message"
	expected := "formatted string"
	suite.Formatter.On("Msg", msg, mock.Anything).Return(expected)
	Msg(msg)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

func (suite *PrinterSuite) TestSMessage() {
	msg := "test message"
	expected := "formatted string"
	suite.Formatter.On("Msg", msg, mock.Anything).Return(expected)
	actual := SMsg(msg)
	suite.Equal(expected, actual)
}

func (suite *PrinterSuite) TestMessageWithArgument() {
	msg := "test message"
	args := "test args"
	expected := "formatted string"
	suite.Formatter.On("Msg", mock.Anything, mock.Anything).Return(expected)
	Msg(msg, args)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

func (suite *PrinterSuite) TestSMessageWithArgument() {
	msg := "test message"
	args := "test args"
	expected := "formatted string"
	suite.Formatter.On("Msg", mock.Anything, mock.Anything).Return(expected)
	actual := SMsg(msg, args)
	suite.Equal(expected, actual)
}

func (suite *PrinterSuite) TestError() {
	errMsg := "test error message"
	expected := "formatted error string"
	suite.Formatter.On("Error", errMsg, mock.Anything).Return(expected)
	Error(errMsg)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

func (suite *PrinterSuite) TestSError() {
	errMsg := "test error message"
	expected := "formatted error string"
	suite.Formatter.On("Error", errMsg, mock.Anything).Return(expected)
	actual := SError(errMsg)
	suite.Equal(expected, actual)
}

func (suite *PrinterSuite) TestErrorWithArgument() {
	errMsg := "test message %v"
	args := "test arg"
	expected := "formatted error string"
	suite.Formatter.On("Error", errMsg, mock.Anything).Return(expected)
	Error(errMsg, args)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

func (suite *PrinterSuite) TestSErrorWithArgument() {
	errMsg := "test message %v"
	args := "test arg"
	expected := "formatted error string"
	suite.Formatter.On("Error", errMsg, mock.Anything).Return(expected)
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
		"row1",
		"row2",
		"row3",
	}
	suite.Formatter.On("Tabulate", table).Return(expected)
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
		"row1",
		"row2",
		"row3",
	}
	suite.Formatter.On("Tabulate", table).Return(expected)
	actual := STabulate(table)
	suite.Equal(expected, actual)
}
func TestPrinterSuite(t *testing.T) {
	suite.Run(t, new(PrinterSuite))
}
