package printer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockWriter is a mock io.Writer for testing
type MockWriter struct {
	mock.Mock
}

func (m *MockWriter) Write(p []byte) (n int, err error) {
	return 0, nil
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

func newMockWriter() *MockWriter {
	writer := new(MockWriter)
	// writer.On("Write", mock.Anything).Return(0, nil)
	return writer
}

func newMockFormatter() *MockFormatter {
	formatter := new(MockFormatter)
	// formatter.On("Msg", mock.Anything, mock.Anything).Return("formatted string")
	// formatter.On("Error", mock.Anything).Return(0, nil)
	// formatter.On("Tabulate", mock.Anything).Return(0, nil)
	return formatter
}

func newMockStenciller() *MockStenciller {
	stenciller := new(MockStenciller)
	// stenciller.On("AddStencil", mock.Anything)
	// stenciller.On("UseStencil", mock.Anything).Return(0, nil)
	return stenciller
}

func (suite *PrinterSuite) SetupTest() {
	suite.Writer = newMockWriter()
	suite.Formatter = newMockFormatter()
	suite.Stenciller = newMockStenciller()
	suite.Printer = &Printer{suite.Writer, suite.Formatter, suite.Stenciller}
}

func (suite *PrinterSuite) TestMessage() {
	msg := "test message"
	// expected := fmt.Sprintf("%v\n", msg)
	suite.Formatter.On("Msg", mock.Anything, mock.Anything).Return("formatted string")
	suite.Printer.Msg(msg)
	// suite.Formatter.AssertExpectations(suite.T())
	suite.Writer.AssertCalled(suite.T(), "Write")
}

func (suite *PrinterSuite) TestMessageWithArgument() {
	msg := "test message %v"
	arg := "arg"
	// expected := fmt.Sprintf("%v\n", fmt.Sprintf(msg, arg))
	suite.Formatter.On("Msg", mock.Anything, mock.Anything).Return("formatted string")
	suite.Printer.Msg(msg, arg)
	// suite.Formatter
	// suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

func (suite *PrinterSuite) TestError() {
	errMsg := "test error message"
	expected := fmt.Sprintf("Error: %v\n", errMsg)
	suite.Printer.Error(errMsg)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

func (suite *PrinterSuite) TestErrorWithArgument() {
	errMsg := "test message %v"
	arg := "arg"
	expected := fmt.Sprintf("Error: %v\n", fmt.Sprintf(errMsg, arg))
	suite.Printer.Error(errMsg, arg)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
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
	suite.Printer.Tabulate(table)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

func TestPrinterSuite(t *testing.T) {
	suite.Run(t, new(PrinterSuite))
}
