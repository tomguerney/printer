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

func (m *MockStenciller) AddStencil(id, template string, colors map[string]string) {
	m.Called(id, template, colors)
}

func (m *MockStenciller) UseStencil(id string, s interface{}) (string, error) {
	args := m.Called(id, s)
	return args.String(0), args.Error(1)
}

type PrinterSuite struct {
	suite.Suite
	Printer    *Printer
	Writer     *MockWriter
	Formatter  *MockFormatter
	Stenciller *MockStenciller
}

func (suite *PrinterSuite) SetupTest() {
	suite.Writer = new(MockWriter)
	suite.Writer.On("Write", mock.Anything).Return(0, nil)
	suite.Formatter = new(MockFormatter)
	suite.Stenciller = new(MockStenciller)
	suite.Printer = &Printer{suite.Writer, suite.Formatter, suite.Stenciller}
}

func (suite *PrinterSuite) TestMessage() {
	msg := "test message"
	formatted := "formatted string"
	suite.Formatter.On("Msg", msg, mock.Anything).Return(formatted)
	suite.Printer.Msg(msg)
	suite.Writer.AssertCalled(suite.T(), "Write", formatted)
}

func (suite *PrinterSuite) TestMessageWithArgument() {
	msg := "test message"
	args := "test args"
	formatted := "formatted string"
	suite.Formatter.On("Msg", mock.Anything, mock.Anything).Return(formatted)
	suite.Printer.Msg(msg, args)
	suite.Writer.AssertCalled(suite.T(), "Write", formatted)
}

func (suite *PrinterSuite) TestError() {
	errMsg := "test error message"
	formatted := "formatted error string"
	suite.Formatter.On("Error", errMsg, mock.Anything).Return(formatted)
	suite.Printer.Error(errMsg)
	suite.Writer.AssertCalled(suite.T(), "Write", formatted)
}

func (suite *PrinterSuite) TestErrorWithArgument() {
	errMsg := "test message %v"
	args := "test arg"
	formatted := "formatted error string"
	suite.Formatter.On("Error", errMsg, mock.Anything).Return(formatted)
	suite.Printer.Error(errMsg, args)
	suite.Writer.AssertCalled(suite.T(), "Write", formatted)
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
	suite.Formatter.On("Tabulate", table).Return(expected)
	suite.Printer.Tabulate(table)
	suite.Writer.AssertCalled(suite.T(), "Write", expected[0])
	suite.Writer.AssertCalled(suite.T(), "Write", expected[1])
	suite.Writer.AssertCalled(suite.T(), "Write", expected[2])
	suite.Writer.AssertNumberOfCalls(suite.T(), "Write", 3)
}

func TestPrinterSuite(t *testing.T) {
	suite.Run(t, new(PrinterSuite))
}
